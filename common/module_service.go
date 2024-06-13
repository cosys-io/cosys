package common

import "fmt"

type ModuleService struct {
	Cosys *Cosys
}

func (e *ModuleService) FindOne(uid string, id int, params MSParams) (Entity, error) {
	queryParams := transformParams(params)

	model, ok := e.Cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	queryParams = queryParams.Where(model.Id_().Eq(id))
	queryParams = queryParams.Limit(1)

	return e.Cosys.Database().FindOne(uid, queryParams)
}

func (e *ModuleService) FindMany(uid string, params MSParams) ([]Entity, error) {
	queryParams := transformParams(params)

	return e.Cosys.Database().FindMany(uid, queryParams)
}

func (e *ModuleService) Create(uid string, entity Entity, params MSParams) (Entity, error) {
	queryParams := transformParams(params)

	return e.Cosys.Database().Create(uid, entity, queryParams)
}

func (e *ModuleService) Update(uid string, entity Entity, id int, params MSParams) (Entity, error) {
	queryParams := transformParams(params)

	model, ok := e.Cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	queryParams = queryParams.Where(model.Id_().Eq(id))

	return e.Cosys.Database().Update(uid, entity, queryParams)
}

func (e *ModuleService) Delete(uid string, id int, params MSParams) (Entity, error) {
	queryParams := transformParams(params)

	model, ok := e.Cosys.Models[uid]
	if !ok {
		return nil, fmt.Errorf("model not found: %s", uid)
	}

	queryParams = queryParams.Where(model.Id_().Eq(id))

	return e.Cosys.Database().Delete(uid, queryParams)
}

type MSParams struct {
	GetFields []Attribute
	SetFields []Attribute
	Filters   []Condition
	StartVal  int64
	LimitVal  int64
	Sorts     []*Order
	Populates []Attribute
}

func MSParam() MSParams {
	return MSParams{
		[]Attribute{},
		[]Attribute{},
		[]Condition{},
		0,
		-1,
		[]*Order{},
		[]Attribute{},
	}
}

func (p MSParams) GetField(fields ...Attribute) MSParams {
	p.GetFields = append(p.GetFields, fields...)
	return p
}

func (p MSParams) SetField(fields ...Attribute) MSParams {
	p.SetFields = append(p.SetFields, fields...)
	return p
}

func (p MSParams) Filter(filters ...Condition) MSParams {
	p.Filters = append(p.Filters, filters...)
	return p
}

func (p MSParams) Limit(limit int64) MSParams {
	p.LimitVal = limit
	return p
}

func (p MSParams) Start(start int64) MSParams {
	p.StartVal = start
	return p
}

func (p MSParams) Sort(sorts ...*Order) MSParams {
	p.Sorts = append(p.Sorts, sorts...)
	return p
}

func (p MSParams) Populate(populates ...Attribute) MSParams {
	p.Populates = append(p.Populates, populates...)
	return p
}

func transformParams(params MSParams) DBParams {
	return DBParam().
		Select(params.GetFields...).
		Insert(params.SetFields...).
		Where(params.Filters...).
		Limit(params.LimitVal).
		Offset(params.StartVal).
		OrderBy(params.Sorts...).
		Populate(params.Populates...)
}
