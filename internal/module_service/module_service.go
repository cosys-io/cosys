package module_service

import (
	"github.com/cosys-io/cosys/internal/common"
	"github.com/cosys-io/cosys/internal/models"
)

type ModuleService struct {
	cosys common.Cosys
}

func NewEntityService(cosys common.Cosys) *ModuleService {
	return &ModuleService{
		cosys,
	}
}

func (e *ModuleService) FindOne(uid string, id int, params *common.ESParams) (models.Entity, error) {
	queryParams := TransformParams(params)

	model, err := e.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	queryParams.Where(model.Id_().Eq(id))
	queryParams.Limit(1)

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.FindOne(queryParams)
}

func (e *ModuleService) FindMany(uid string, params *common.ESParams) ([]models.Entity, error) {
	queryParams := TransformParams(params)

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.FindMany(queryParams)
}

func (e *ModuleService) Create(uid string, entity models.Entity, params *common.ESParams) (models.Entity, error) {
	queryParams := TransformParams(params)

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.Create(entity, queryParams)
}

func (e *ModuleService) Update(uid string, entity models.Entity, id int, params *common.ESParams) (models.Entity, error) {
	queryParams := TransformParams(params)

	model, err := e.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	queryParams.Where(model.Id_().Eq(id))

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.Update(entity, queryParams)
}

func (e *ModuleService) Delete(uid string, id int, params *common.ESParams) (models.Entity, error) {
	queryParams := TransformParams(params)

	model, err := e.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	queryParams.Where(model.Id_().Eq(id))

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.Delete(queryParams)
}
