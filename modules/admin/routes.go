package admin

import (
	"encoding/json"
	"fmt"
	"github.com/cosys-io/cosys/common"
	"net/http"
	"strconv"
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
		"findOne": findOne(uid, model),
		"create":  create(uid, model),
		"update":  update(uid, model),
		"delete":  delete(uid, model),
	}
	adminModule.Controllers[modelName+"Admin"] = controller

	routes := []*common.Route{
		common.NewRoute("GET", fmt.Sprintf(`/admin/%s/{documentId}`, modelName), modelName+"Admin.findOne"),
		common.NewRoute("POST", fmt.Sprintf(`/admin/%s`, modelName), modelName+"Admin.create"),
		common.NewRoute("PUT", fmt.Sprintf(`/admin/%s/{documentId}`, modelName), modelName+"Admin.update"),
		common.NewRoute("DELETE", fmt.Sprintf(`/admin/%s/{documentId}`, modelName), modelName+"Admin.delete"),
	}
	adminModule.Routes = append(adminModule.Routes, routes...)

	return nil
}

func findOne(uid string, contentModel common.Model) func(common.Cosys) http.HandlerFunc {
	return func(cs common.Cosys) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			params, err := common.ReadParams(r)
			if err != nil {
				common.RespondInternalError(w)
				return
			}

			id, err := strconv.Atoi(params["documentId"])
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

			id, err := strconv.Atoi(params["documentId"])
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

			id, err := strconv.Atoi(params["documentId"])
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
