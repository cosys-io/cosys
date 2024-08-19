package routes

import (
	"encoding/json"
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/server/response"
	"net/http"
)

// FindMany returns the find many ActionFunc for the model of the given uid.
func FindMany(modelUid string) common.ActionFunc {
	return func(cosys *common.Cosys) (http.HandlerFunc, error) {
		model, err := cosys.Model(modelUid)
		if err != nil {
			return nil, err
		}

		database, err := cosys.Database()
		if err != nil {
			return nil, err
		}

		return func(w http.ResponseWriter, r *http.Request) {
			page, err := getPage(r)
			if err != nil {
				response.RespondError(w, "Could not find "+model.PluralHumanName_(), http.StatusBadRequest)
				return
			}

			dbParams, err := getParams(r, model.Attributes_())
			if err != nil {
				response.RespondError(w, "Could not find "+model.PluralHumanName_(), http.StatusBadRequest)
				return
			}

			entities, err := database.FindMany(modelUid, dbParams)
			if err != nil {
				response.RespondError(w, "Could not find "+model.PluralHumanName_(), http.StatusBadRequest)
				return
			}

			response.RespondMany(w, entities, page, http.StatusOK)
		}, nil
	}
}

// FindOne returns the find one ActionFunc for the model of the given uid.
func FindOne(modelUid string) common.ActionFunc {
	return func(cosys *common.Cosys) (http.HandlerFunc, error) {
		model, err := cosys.Model(modelUid)
		if err != nil {
			return nil, err
		}

		database, err := cosys.Database()
		if err != nil {
			return nil, err
		}

		return func(w http.ResponseWriter, r *http.Request) {
			id, err := getId(r)
			if err != nil {
				response.RespondError(w, "Could not find "+model.SingularHumanName_(), http.StatusBadRequest)
				return
			}

			dbParams := common.NewDBParamsBuilder().
				Where(model.IdAttribute_().(common.IntAttribute).Eq(id)).
				Build()

			entity, err := database.FindOne(modelUid, dbParams)
			if err != nil {
				response.RespondError(w, "Could not find "+model.SingularHumanName_(), http.StatusBadRequest)
				return
			}

			response.RespondOne(w, entity, http.StatusOK)
		}, nil
	}
}

// Create returns the create ActionFunc for the model of the given uid.
func Create(modelUid string) common.ActionFunc {
	return func(cosys *common.Cosys) (http.HandlerFunc, error) {
		model, err := cosys.Model(modelUid)
		if err != nil {
			return nil, err
		}

		database, err := cosys.Database()
		if err != nil {
			return nil, err
		}

		return func(w http.ResponseWriter, r *http.Request) {
			entity := model.New_()

			if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
				response.RespondError(w, "Could not create "+model.SingularHumanName_(), http.StatusBadRequest)
				return
			}

			newEntity, err := database.Create(modelUid, entity, common.NewDBParams())
			if err != nil {
				response.RespondError(w, "Could not create "+model.SingularHumanName_(), http.StatusBadRequest)
				return
			}

			response.RespondOne(w, newEntity, http.StatusOK)
		}, nil
	}
}

// Update returns the update ActionFunc for the model of the given uid.
func Update(modelUid string) common.ActionFunc {
	return func(cosys *common.Cosys) (http.HandlerFunc, error) {
		model, err := cosys.Model(modelUid)
		if err != nil {
			return nil, err
		}

		database, err := cosys.Database()
		if err != nil {
			return nil, err
		}

		return func(w http.ResponseWriter, r *http.Request) {
			id, err := getId(r)
			if err != nil {
				response.RespondError(w, "Could not update "+model.SingularHumanName_(), http.StatusBadRequest)
				return
			}

			entity := model.New_()

			if err := json.NewDecoder(r.Body).Decode(entity); err != nil {
				response.RespondError(w, "Could not update "+model.SingularHumanName_(), http.StatusBadRequest)
				return
			}

			dbParams := common.NewDBParamsBuilder().
				Where(model.IdAttribute_().(common.IntAttribute).Eq(id)).
				Build()

			newEntity, err := database.Update(modelUid, entity, dbParams)
			if err != nil {
				response.RespondError(w, "Could not update "+model.SingularHumanName_(), http.StatusBadRequest)
				return
			}

			response.RespondOne(w, newEntity, http.StatusOK)
		}, nil
	}
}

// Delete returns the delete ActionFunc for the model of the given uid.
func Delete(modelUid string) common.ActionFunc {
	return func(cosys *common.Cosys) (http.HandlerFunc, error) {
		model, err := cosys.Model(modelUid)
		if err != nil {
			return nil, err
		}

		database, err := cosys.Database()
		if err != nil {
			return nil, err
		}

		return func(w http.ResponseWriter, r *http.Request) {
			id, err := getId(r)
			if err != nil {
				response.RespondError(w, "Could not delete "+model.SingularHumanName_(), http.StatusBadRequest)
				return
			}

			dbParams := common.NewDBParamsBuilder().
				Where(model.IdAttribute_().(common.IntAttribute).Eq(id)).
				Build()

			entity, err := database.Delete(modelUid, dbParams)
			if err != nil {
				response.RespondError(w, "Could not delete "+model.SingularHumanName_(), http.StatusBadRequest)
				return
			}

			response.RespondOne(w, entity, http.StatusOK)
		}, nil
	}
}
