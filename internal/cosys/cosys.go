package cosys

import (
	"database/sql"

	"github.com/cosys-io/cosys/internal/common"
	genlifecycles "github.com/cosys-io/cosys/internal/gen/lifecycles"
	genmodels "github.com/cosys-io/cosys/internal/gen/models"
	genservices "github.com/cosys-io/cosys/internal/gen/services"
	"github.com/cosys-io/cosys/internal/models"
	"github.com/cosys-io/cosys/internal/module_service"
	"github.com/cosys-io/cosys/internal/query_engine"
)

type Cosys struct {
	DB *sql.DB
}

func (c *Cosys) QueryEngine(uid string) (common.QueryEngine, error) {
	return query_engine.NewQueryEngine(c, c.DB, "sqlite3", uid)
}

func (c *Cosys) ModuleService() (common.EntityService, error) {
	return module_service.NewEntityService(c), nil
}

func (c *Cosys) Service(uid string) (common.ServiceFunction, error) {
	return genservices.Service(c, uid)
}

func (c *Cosys) Model(uid string) (models.Model, error) {
	return genmodels.Model(uid)
}

func (c *Cosys) Lifecycle(uid string) (*common.Lifecycles, error) {
	return genlifecycles.Lifecycle(uid)
}
