package content_builder

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/cosys_cli/cmd"
	"gopkg.in/yaml.v3"
	"net/http"
)

var Controller = map[string]common.Action{
	"schema": schema,
	"build":  build,
}

func schema(cosys common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schemas := []common.ModelSchema{}

		for _, model := range cosys.Models {
			schemas = append(schemas, *model.Schema_())
		}

		common.RespondMany(w, schemas, 1, http.StatusOK)
	}
}

func build(cosys common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
		}

		if len(params) == 0 {
			common.RespondInternalError(w)
		}

		name := params[0]

		schemaParsed := &common.ModelSchemaParsed{}

		if err := yaml.NewDecoder(r.Body).Decode(schemaParsed); err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		if err := cmd.GenerateType(name, schemaParsed.Schema()); err != nil {
			common.RespondError(w, "Unable to build content type.", http.StatusBadRequest)
		}

		common.RespondOne(w, "Content type successfully created.", 200)
	}
}
