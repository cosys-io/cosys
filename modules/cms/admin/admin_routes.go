package admin

import (
	"github.com/cosys-io/cosys/common"
	"github.com/cosys-io/cosys/modules/cms/routes"
)

func AddAdminRoutes(cosys *common.Cosys, models map[string]common.Model) error {
	adminRoutes := make([]common.Route, len(models)*5)

	index := 0
	for modelUid, model := range models {
		modelApi := model.PluralKebabName_()

		adminRoutes[index] = common.NewRoute("GET", `/admin/`+modelApi, routes.FindMany(modelUid))
		adminRoutes[index+1] = common.NewRoute("GET", `/admin/`+modelApi+`/{id}`, routes.FindOne(modelUid))
		adminRoutes[index+2] = common.NewRoute("POST", `/admin/`+modelApi, routes.Create(modelUid))
		adminRoutes[index+3] = common.NewRoute("PUT", `/admin/`+modelApi+`/{id}`, routes.Update(modelUid))
		adminRoutes[index+4] = common.NewRoute("DELETE", `/admin/`+modelApi+`/{id}`, routes.Delete(modelUid))
		index += 5
	}

	return cosys.AddRoutes(adminRoutes...)
}
