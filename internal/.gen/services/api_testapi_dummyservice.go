// AUTO-GENERATED
// DO NOT CHANGE

package services

import (
	. "github.com/cosys-io/cosys/internal/.gen/models"
	"github.com/cosys-io/cosys/internal/cosys"

	"github.com/cosys-io/cosys/internal/api/testapi/services"
)

type Api_Testapi_Dummyservice struct {
	cs cosys.Cosys
}

func New_Api_Testapi_Dummyservice(cs cosys.Cosys) any {
	return &Api_Testapi_Dummyservice{
		cs,
	}
}

func (a *Api_Testapi_Dummyservice) FindUserByName(name string) (*User, error) {
	return services.Services["dummy"].Function("findUserByName")(a.cs).(func(name string) (*User, error))(name)
}

func (a *Api_Testapi_Dummyservice) FindActiveUsers() ([]*User, error) {
	return services.Services["dummy"].Function("findActiveUsers")(a.cs).(func() ([]*User, error))()
}
