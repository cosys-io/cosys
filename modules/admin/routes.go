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
	for modelUid, model := range cosys.Models {
		if err := AddRoutes(modelUid, model, &cosys); err != nil {
			return common.Cosys{}, err
		}
	}

	return cosys, nil
}

func AddRoutes(modelUid string, model common.Model, cosys *common.Cosys) error {
	if cosys == nil {
		return fmt.Errorf("cosys is nil")
	}

	adminModule, ok := cosys.Modules["admin"]
	if !ok {
		return fmt.Errorf("admin module not found")
	}

	modelName := model.Name_()

	controller := map[string]common.Action{
		"findMany": findManyEntity(modelUid, model.Schema_().PluralName),
		"findOne":  findOneEntity(modelUid, model.Schema_().SingularName),
		"create":   createEntity(modelUid, model.Schema_().SingularName),
		"update":   updateEntity(modelUid, model.Schema_().SingularName),
		"delete":   deleteEntity(modelUid, model.Schema_().SingularName),
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

func findManyEntity(modelUid, modelName string) func(common.Cosys) http.HandlerFunc {
	return func(cosys common.Cosys) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params, err := common.ReadParams(r)
			if err != nil {
				common.RespondInternalError(w)
				return
			}

			model, ok := cosys.Models[modelUid]
			if !ok {
				common.RespondInternalError(w)
				return
			}
			attrs := model.All_()

			page := 1
			pageSize := int64(20)
			var sort []*common.Order
			var filter []common.Condition
			var fields []common.Attribute
			var populate []common.Attribute

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
				sortStrings := strings.Split(sortSliceString, ",")
				for _, sortString := range sortStrings {
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

					for _, attr := range attrs {
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
				fieldStrings := strings.Split(fieldSliceString, ",")
				for _, fieldString := range fieldStrings {
					var fieldAttr common.Attribute

					for _, attr := range attrs {
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
				populateStrings := strings.Split(populateSliceString, ",")
				for _, populateString := range populateStrings {
					var populateAttr common.Attribute

					for _, attr := range attrs {
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

			msParams := common.NewMSParamsBuilder().
				Start(pageSize * (int64(page) - 1)).
				Limit(pageSize).
				Sort(sort...).
				Filter(filter...).
				GetField(fields...).
				Populate(populate...).
				Build()
			entities, err := cosys.ModuleService().FindMany(modelUid, msParams)
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not find %s.", modelName), http.StatusBadRequest)
				return
			}

			common.RespondMany(w, entities, page, http.StatusOK)
		}
	}
}

func findOneEntity(modelUid, modelName string) func(common.Cosys) http.HandlerFunc {
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

			entity, err := cs.ModuleService().FindOne(modelUid, id, common.NewMSParams())
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not find %s.", modelName), http.StatusBadRequest)
				return
			}

			common.RespondOne(w, entity, http.StatusOK)
		}
	}
}

func createEntity(modelUid, modelName string) func(common.Cosys) http.HandlerFunc {
	return func(cs common.Cosys) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			model, ok := cs.Models[modelUid]
			if !ok {
				common.RespondInternalError(w)
				return
			}
			entity := model.New_()

			if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			newEntity, err := cs.ModuleService().Create(modelUid, entity, common.NewMSParams())
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not create %s.", modelName), http.StatusBadRequest)
				return
			}

			common.RespondOne(w, newEntity, http.StatusOK)
		}
	}
}

func updateEntity(modelUid, modelName string) func(common.Cosys) http.HandlerFunc {
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

			model, ok := cs.Models[modelUid]
			if !ok {
				common.RespondInternalError(w)
				return
			}
			entity := model.New_()

			if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
				common.RespondError(w, "Bad request.", http.StatusBadRequest)
				return
			}

			newEntity, err := cs.ModuleService().Update(modelUid, entity, id, common.NewMSParams())
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not update %s.", modelName), http.StatusBadRequest)
				return
			}

			common.RespondOne(w, newEntity, http.StatusOK)
		}
	}
}

func deleteEntity(modelUid, modelName string) func(common.Cosys) http.HandlerFunc {
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

			oldEntity, err := cs.ModuleService().Delete(modelUid, id, common.NewMSParams())
			if err != nil {
				common.RespondError(w, fmt.Sprintf("Could not delete %s.", modelName), http.StatusBadRequest)
				return
			}

			common.RespondOne(w, oldEntity, http.StatusOK)
		}
	}
}
