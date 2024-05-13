package cosys

import "github.com/cosys-io/cosys/internal/models"

type EntityService interface {
	FindOne(string, int, *ESParams) (models.Entity, error)
	FindMany(string, *ESParams) ([]models.Entity, error)
	Create(string, models.Entity, *ESParams) (models.Entity, error)
	Update(string, models.Entity, int, *ESParams) (models.Entity, error)
	Delete(string, int, *ESParams) (models.Entity, error)
}

type ESParams struct {
	GetFields []models.IAttribute
	SetFields []models.IAttribute
	Filters   []models.Condition
	StartVal  int64
	LimitVal  int64
	Sorts     []*models.Order
	Populates []models.IAttribute
}

func ESParam() *ESParams {
	return &ESParams{
		[]models.IAttribute{},
		[]models.IAttribute{},
		[]models.Condition{},
		0,
		-1,
		[]*models.Order{},
		[]models.IAttribute{},
	}
}

func (p *ESParams) GetField(fields ...models.IAttribute) *ESParams {
	p.GetFields = append(p.GetFields, fields...)
	return p
}

func (p *ESParams) SetField(fields ...models.IAttribute) *ESParams {
	p.SetFields = append(p.SetFields, fields...)
	return p
}

func (p *ESParams) Filter(filters ...models.Condition) *ESParams {
	p.Filters = append(p.Filters, filters...)
	return p
}

func (p *ESParams) Limit(limit int64) *ESParams {
	p.LimitVal = limit
	return p
}

func (p *ESParams) Start(start int64) *ESParams {
	p.StartVal = start
	return p
}

func (p *ESParams) Sort(sorts ...*models.Order) *ESParams {
	p.Sorts = append(p.Sorts, sorts...)
	return p
}

func (p *ESParams) Populate(populates ...models.IAttribute) *ESParams {
	p.Populates = append(p.Populates, populates...)
	return p
}
