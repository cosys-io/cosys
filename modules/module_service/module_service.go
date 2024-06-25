package module_service

import (
	"fmt"

	"github.com/cosys-io/cosys/common"
)

type ModuleService struct {
	Cosys *common.Cosys
}

func (e ModuleService) FindOne(uid string, id int, params common.MSParams) (common.Entity, error) {
	queryParams := transformParams(params)

	model, ok := e.Cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	queryParams.Where = append(queryParams.Where, model.Id_().Eq(id))
	queryParams.Limit = 1

	return e.Cosys.Database().FindOne(uid, queryParams)
}

func (e ModuleService) FindMany(uid string, params common.MSParams) ([]common.Entity, error) {
	queryParams := transformParams(params)

	return e.Cosys.Database().FindMany(uid, queryParams)
}

func (e ModuleService) Create(uid string, entity common.Entity, params common.MSParams) (common.Entity, error) {
	queryParams := transformParams(params)

	return e.Cosys.Database().Create(uid, entity, queryParams)
}

func (e ModuleService) Update(uid string, entity common.Entity, id int, params common.MSParams) (common.Entity, error) {
	queryParams := transformParams(params)

	model, ok := e.Cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	queryParams.Where = append(queryParams.Where, model.Id_().Eq(id))

	return e.Cosys.Database().Update(uid, entity, queryParams)
}

func (e ModuleService) Delete(uid string, id int, params common.MSParams) (common.Entity, error) {
	queryParams := transformParams(params)

	model, ok := e.Cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	queryParams.Where = append(queryParams.Where, model.Id_().Eq(id))

	return e.Cosys.Database().Delete(uid, queryParams)
}

func transformParams(params common.MSParams) common.DBParams {
	return common.NewDBParamsBuilder().
		Select(params.Select...).
		Insert(params.Fields...).
		Where(params.Filter...).
		Limit(params.Limit).
		Offset(params.Start).
		OrderBy(params.Sort...).
		Populate(params.Populate...).
		Build()
}
