package entityservice

import (
	"github.com/cosys-io/cosys/internal/cosys"
	"github.com/cosys-io/cosys/internal/models"
)

type EntityService struct {
	cosys cosys.Cosys
}

func NewEntityService(cosys cosys.Cosys) *EntityService {
	return &EntityService{
		cosys,
	}
}

func (e *EntityService) FindOne(uid string, id int, params *cosys.ESParams) (models.Entity, error) {
	queryParams := TransformParams(params)

	model, err := e.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	queryParams.Where(model.Model_Id().Eq(id))
	queryParams.Limit(1)

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.FindOne(queryParams)
}

func (e *EntityService) FindMany(uid string, params *cosys.ESParams) ([]models.Entity, error) {
	queryParams := TransformParams(params)

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.FindMany(queryParams)
}

func (e *EntityService) Create(uid string, entity models.Entity, params *cosys.ESParams) (models.Entity, error) {
	queryParams := TransformParams(params)

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.Create(entity, queryParams)
}

func (e *EntityService) Update(uid string, entity models.Entity, id int, params *cosys.ESParams) (models.Entity, error) {
	queryParams := TransformParams(params)

	model, err := e.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	queryParams.Where(model.Model_Id().Eq(id))

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.Update(entity, queryParams)
}

func (e *EntityService) Delete(uid string, id int, params *cosys.ESParams) (models.Entity, error) {
	queryParams := TransformParams(params)

	model, err := e.cosys.Model(uid)
	if err != nil {
		return nil, err
	}

	queryParams.Where(model.Model_Id().Eq(id))

	qe, err := e.cosys.QueryEngine(uid)
	if err != nil {
		return nil, err
	}

	return qe.Delete(queryParams)
}
