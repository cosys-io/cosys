// AUTO-GENERATED
// DO NOT CHANGE

package models

import (
	"fmt"

	"github.com/cosys-io/cosys/internal/models"
)

var genModels = map[string]models.Model{
	"testapi::users": Users,
}

func Model(uid string) (models.Model, error) {
	model := genModels[uid]
	if model == nil {
		return nil, fmt.Errorf("invalid uid: %s", uid)
	}

	return model, nil
}
