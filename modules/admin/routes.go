package admin

import (
	"encoding/json"
	"fmt"
	"github.com/cosys-io/cosys/common"
	"net/http"
	"strconv"
	"strings"
)

func OnRegister(cosys common.Cosys) (common.Cosys, error) {
	for uid, model := range cosys.Models {
		if err := AddRoutes(uid, model, &cosys); err != nil {
			return cosys, err
		}
	}

	return cosys, nil
}

func AddRoutes(uid string, model common.Model, cosys *common.Cosys) error {
	adminModule, ok := cosys.Modules["admin"]
	if !ok {
		return fmt.Errorf("admin module not found")
	}

	modelName := model.Name_()

	controller := map[string]common.Action{
		"findMany": findMany(uid, model),
		"findOne":  findOne(uid, model),
		"create":   create(uid, model),
		"update":   update(uid, model),
		"delete":   delete(uid, model),
	}
	adminModule.Controllers[modelName+"Admin"] = controller

	routes := []*common.Route{
		common.NewRoute("GET", fmt.Sprintf(`/admin/%s`, modelName), modelName+"Admin.findMany"),
		common.NewRoute("GET", fmt.Sprintf(`/admin/%s/{documentId}`, modelName), modelName+"Admin.findOne"),
		common.NewRoute("POST", fmt.Sprintf(`/admin/%s`, modelName), modelName+"Admin.create"),
		common.NewRoute("PUT", fmt.Sprintf(`/admin/%s/{documentId}`, modelName), modelName+"Admin.update"),
		common.NewRoute("DELETE", fmt.Sprintf(`/admin/%s/{documentId}`, modelName), modelName+"Admin.delete"),
	}
	adminModule.Routes = append(adminModule.Routes, routes...)

	return nil
}

func findMany(uid string, contentModel common.Model) func(common.Cosys) http.HandlerFunc {
	return func(cosys common.Cosys) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params, err := common.ReadParams(r)
			if err != nil {
				common.RespondInternalError(w)
				return
			}

			model, ok := cosys.Models[uid]
			if !ok {
				common.RespondInternalError(w)
				return
			}
			attrSlice := model.All_()

			page := 1
			pageSize := int64(20)
			sort := []*common.Order{}
			filter := []common.Condition{}
			fields := []common.Attribute{}
			populate := []common.Attribute{}

			pageSizeString, ok := params["pageSize"]
			if ok {
				pageSize, err = strconv.ParseInt(pageSizeString, 10, 64)
				if err != nil {
					common.RespondError(w, "Bad request.", http.StatusBadRequest)
					return
				}
			}

			pageString, ok := params["page"]
			if ok {
				page, err = strconv.Atoi(pageString)
				if err != nil {
					common.RespondError(w, "Bad request.", http.StatusBadRequest)
					return
				}
			}

			sortSliceString, ok := params["sort"]
			if ok {
				sortSlice := strings.Split(sortSliceString, ",")
				for _, sortString := range sortSlice {
					if len(sortString) == 0 {
						common.RespondError(w, "Bad request.", http.StatusBadRequest)
						return
					}

					isAsc := true
					if sortString[0] == '-' {
						isAsc = false
						sortString = sortString[1:]
					}

					var sortAttr common.Attribute

					for _, attr := range attrSlice {
						if attr.Name() == sortString {
							sortAttr = attr
						}
					}

					if sortAttr == nil {
						common.RespondError(w, "Bad request.", http.StatusBadRequest)
						return
					}

					if isAsc {
						sort = append(sort, sortAttr.Asc())
					} else {
						sort = append(sort, sortAttr.Desc())
					}
				}
			}

			fieldSliceString, ok := params["fields"]
			if ok {
				fieldSlice := strings.Split(fieldSliceString, ",")
				for _, fieldString := range fieldSlice {
					var fieldAttr common.Attribute

					for _, attr := range attrSlice {
						if attr.Name() == fieldString {
							fieldAttr = attr
						}
					}

					if fieldAttr == nil {
						common.RespondError(w, "Bad request.", http.StatusBadRequest)
						return
					}

					fields = append(fields, fieldAttr)
				}
			}

			populateSliceString, ok := params["populate"]
			if ok {
				populateSlice := strings.Split(populateSliceString, ",")
				for _, populateString := range populateSlice {
					var populateAttr common.Attribute

					for _, attr := range attrSlice {
						if attr.Name() == populateString {
							populateAttr = attr
						}
					}

					if populateAttr == nil {
						common.RespondError(w, "Bad request.", http.StatusBadRequest)
						return
					}

					populate = append(populate, populateAttr)
				}
			}

			msParams := common.MSParam().
				Start(pageSize * (int64(page) - 1)).
				Limit(pageSize).
				Sort(sort...).
				Filter(filter...).
				GetField(fields...).
				Populate(populate...)
			entities, err := cosys.ModuleService().FindMany(uid, msParams)
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not find %s.", contentModel.Schema_().PluralName), http.StatusBadRequest)
				return
			}

			common.RespondMany(w, entities, page, http.StatusOK)
		}
	}
}

func findOne(uid string, contentModel common.Model) func(common.Cosys) http.HandlerFunc {
	return func(cs common.Cosys) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params, err := common.ReadParams(r)
			if err != nil {
				common.RespondInternalError(w)
				return
			}

			idString, ok := params["documentId"]
			if !ok {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(idString)
			if err != nil {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			entity, err := cs.ModuleService().FindOne(uid, id, common.MSParam())
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not find %s.", contentModel.Schema_().SingularName), http.StatusBadRequest)
				return
			}

			common.RespondOne(w, entity, http.StatusOK)
		}
	}
}

func create(uid string, contentModel common.Model) func(common.Cosys) http.HandlerFunc {
	return func(cs common.Cosys) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			model, ok := cs.Models[uid]
			if !ok {
				common.RespondInternalError(w)
				return
			}
			entity := model.New_()

			if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			newEntity, err := cs.ModuleService().Create(uid, entity, common.MSParam())
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not create %s.", contentModel.Schema_().SingularName), http.StatusBadRequest)
				return
			}

			common.RespondOne(w, newEntity, http.StatusOK)
		}
	}
}

func update(uid string, contentModel common.Model) func(common.Cosys) http.HandlerFunc {
	return func(cs common.Cosys) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params, err := common.ReadParams(r)
			if err != nil {
				common.RespondInternalError(w)
				return
			}

			idString, ok := params["documentId"]
			if !ok {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(idString)
			if err != nil {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			model, ok := cs.Models[uid]
			if !ok {
				common.RespondInternalError(w)
				return
			}
			entity := model.New_()

			if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			newEntity, err := cs.ModuleService().Update(uid, entity, id, common.MSParam())
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not update %s.", contentModel.Schema_().SingularName), http.StatusBadRequest)
				return
			}

			common.RespondOne(w, newEntity, http.StatusOK)
		}
	}
}

func delete(uid string, contentModel common.Model) func(common.Cosys) http.HandlerFunc {
	return func(cs common.Cosys) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params, err := common.ReadParams(r)
			if err != nil {
				common.RespondInternalError(w)
				return
			}

			idString, ok := params["documentId"]
			if !ok {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			id, err := strconv.Atoi(idString)
			if err != nil {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			oldEntity, err := cs.ModuleService().Delete(uid, id, common.MSParam())
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not delete %s.", contentModel.Schema_().SingularName), http.StatusBadRequest)
				return
			}

			common.RespondOne(w, oldEntity, http.StatusOK)
		}
	}
}
