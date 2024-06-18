package content_builder

import (
	"github.com/cosys-io/cosys/common"
	"gopkg.in/yaml.v3"
	"net/http"
)

var Controller = map[string]common.Action{
	"get":   get,
	"build": build,
}

func get(cosys common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := common.ReadParams(r)
		if err != nil {
			common.RespondInternalError(w)
		}

		if len(params) == 0 {
			common.RespondInternalError(w)
		}

		uid := params[0]
		model, ok := cosys.Models[uid]
		if !ok {
			common.RespondError(w, "Content type not found.", http.StatusNotFound)
		}

		schema := model.Schema_()
		common.RespondOne(w, schema, 200)
	}
}

func build(cosys common.Cosys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schema := common.ModelSchema{}

		if err := yaml.NewDecoder(r.Body).Decode(schema); err != nil {
			common.RespondError(w, "Bad request.", http.StatusBadRequest)
			return
		}

		// TODO: Generate code

		common.RespondOne(w, "Content type successfully created.", 200)
	}
}
