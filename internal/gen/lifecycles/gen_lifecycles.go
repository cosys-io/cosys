package lifecycles

import (
	"fmt"

	api_testapi_users "github.com/cosys-io/cosys/internal/api/testapi/users/lifecycles"
	"github.com/cosys-io/cosys/internal/common"
)

var lifecycles = map[string]*common.Lifecycles{
	"api::testapi.users": api_testapi_users.Lifecycles,
}

func Lifecycle(uid string) (*common.Lifecycles, error) {
	lifecycles := lifecycles[uid]
	if lifecycles == nil {
		return nil, fmt.Errorf("invalid uid: %s", uid)
	}

	return lifecycles, nil
}
