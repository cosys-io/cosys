package common

import "github.com/cosys-io/cosys/internal/models"

type Cosys interface {
	QueryEngine(uid string) (QueryEngine, error)
	ModuleService() (EntityService, error)
	Service(uid string) (ServiceFunction, error)

	Model(uid string) (models.Model, error)
	Lifecycle(uid string) (*Lifecycles, error)
}
