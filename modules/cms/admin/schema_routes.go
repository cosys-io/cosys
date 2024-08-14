package admin

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/cms/generators"
	"github.com/cosys-io/cosys/modules/cms/schema"
	"github.com/cosys-io/cosys/modules/server/response"
	"net/http"
)

// AddSchemaRoutes registers routes for getting and creating schemas for the given models.
func AddSchemaRoutes(cosys *common.Cosys, models map[string]common.Model) error {
	schemaRoutes := []common.Route{
		common.NewRoute("GET", `/admin/schema`, getSchema(models)),
		common.NewRoute("POST", `/admin/schema`, createSchema),
	}

	return cosys.AddRoutes(schemaRoutes...)
}

// getSchema returns the ActionFunc for getting schemas.
func getSchema(models map[string]common.Model) common.ActionFunc {
	schemas := make([]common.ModelSchema, len(models))

	index := 0
	for _, model := range models {
		schemas[index] = model.Schema_()
		index = index + 1
	}

	return func(cosys *common.Cosys) (http.HandlerFunc, error) {
		return func(w http.ResponseWriter, r *http.Request) {
			response.RespondMany(w, schemas, 1, http.StatusOK)
		}, nil
	}
}

// createSchema is the ActionFunc for creating schemas.
var createSchema common.ActionFunc = func(cosys *common.Cosys) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		var newSchema schema.ModelSchema
		if err := schema.ParseSchema(&newSchema, r.Body); err != nil {
			response.RespondError(w, "Could not create content type.", http.StatusBadRequest)
			return
		}

		if err := generators.GenerateType(&newSchema); err != nil {
			response.RespondError(w, "Could not create content type.", http.StatusBadRequest)
			return
		}

		response.RespondOne(w, nil, http.StatusOK)
	}, nil
}
