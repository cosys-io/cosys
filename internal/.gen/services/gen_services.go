// AUTO-GENERATED
// DO NOT CHANGE

package services

import (
	"fmt"

	"github.com/cosys-io/cosys/internal/cosys"
)

var genServices = map[string](func(cosys.Cosys) any){
	"api::testapi.dummyservice": New_Api_Testapi_Dummyservice,
}

func Service(cs cosys.Cosys, uid string) (any, error) {
	service := genServices[uid]
	if service == nil {
		return nil, fmt.Errorf("invalid uid: %s", uid)
	}

	return service(cs), nil
}
