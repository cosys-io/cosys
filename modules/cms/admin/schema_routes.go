package admin

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/cms/generators"
	"github.com/cosys-io/cosys/modules/cms/schema"
	"net/http"
)

func AddSchemaRoutes(cosys *common.Cosys, models map[string]common.Model) error {
	schemaRoutes := []common.Route{
		common.NewRoute("GET", `/admin/schema`, getSchema(models)),
		common.NewRoute("POST", `/admin/schema`, createSchema),
	}

	return cosys.AddRoutes(schemaRoutes...)
}

func getSchema(models map[string]common.Model) common.ActionFunc {
	schemas := make([]common.ModelSchema, len(models))

	index := 0
	for _, model := range models {
		schemas[index] = model.Schema_()
		index = index + 1
	}

	return func(cosys *common.Cosys) (http.HandlerFunc, error) {
		return func(w http.ResponseWriter, r *http.Request) {
			// TODO: Respond schemas
		}, nil
	}
}

var createSchema common.ActionFunc = func(cosys *common.Cosys) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		var newSchema schema.ModelSchema
		if err := schema.ParseSchema(&newSchema, r.Body); err != nil {
			// TODO: Bad Request
			return
		}

		if err := generators.GenerateType(&newSchema); err != nil {
			// TODO: Bad Request
			return
		}

		// TODO: Respond success
	}, nil
}
