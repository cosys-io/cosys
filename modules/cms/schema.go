package cms

import (
	"encoding/json"
	"github.com/cosys-io/cosys/common"
	"net/http"
)

var schemaRoutes = []*common.Route{
	common.NewRoute("GET", `/admin/schema`, "admin.schema"),
	common.NewRoute("POST", `/admin/schema`, "admin.build"),
}

var schemaController = map[string]common.Action{
	"schema": schema,
	"build":  build,
}

func schema(cosys common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var schemas []common.ModelSchema

		for _, model := range cosys.Models {
			schemas = append(schemas, *model.Schema_())
		}

		common.RespondMany(w, schemas, 1, http.StatusOK)
	}
}

func build(cosys common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schemaParsed := common.ModelSchemaParsed{}

		if err := json.NewDecoder(r.Body).Decode(&schemaParsed); err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		newSchema, err := schemaParsed.Schema()
		if err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}
		newSchema.Attributes = append([]*common.AttributeSchema{&common.IdSchema}, newSchema.Attributes...)

		if err = generateType(newSchema); err != nil {
			common.RespondError(w, "Unable to build content type.", http.StatusBadRequest)
			return
		}

		common.RespondOne(w, "Content type successfully created.", http.StatusOK)
	}
}
