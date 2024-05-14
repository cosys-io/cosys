// AUTO-GENERATED
// DO NOT CHANGE

package apis

import (
	"github.com/cosys-io/cosys/internal/apis"

	testapi_controllers "github.com/cosys-io/cosys/internal/api/testapi/controllers"
	testapi_middlewares "github.com/cosys-io/cosys/internal/api/testapi/middlewares"
	testapi_policies "github.com/cosys-io/cosys/internal/api/testapi/policies"
	testapi_routes "github.com/cosys-io/cosys/internal/api/testapi/routes"
)

var Apis = map[string]*apis.API{
	"api::testapi": apis.NewAPI(
		testapi_routes.Routes,
		testapi_controllers.Controllers,
		testapi_middlewares.Middlewares,
		testapi_policies.Policies,
	),
}
