package cosys

import (
	"database/sql"

	"github.com/cosys-io/cosys/internal/common"
	"github.com/cosys-io/cosys/internal/entityservice"
	genlifecycles "github.com/cosys-io/cosys/internal/gen/lifecycles"
	genmodels "github.com/cosys-io/cosys/internal/gen/models"
	genservices "github.com/cosys-io/cosys/internal/gen/services"
	"github.com/cosys-io/cosys/internal/models"
	"github.com/cosys-io/cosys/internal/queryengine"
)

type Cosys struct {
	DB *sql.DB
}

func (c *Cosys) QueryEngine(uid string) (common.QueryEngine, error) {
	return queryengine.NewQueryEngine(c, c.DB, "sqlite3", uid)
}

func (c *Cosys) EntityService() (common.EntityService, error) {
	return entityservice.NewEntityService(c), nil
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
