package cosys

import "github.com/cosys-io/cosys/internal/models"

type QueryEngine interface {
	FindOne(*QEParams) (models.Entity, error)
	FindMany(*QEParams) ([]models.Entity, error)
	Create(models.Entity, *QEParams) (models.Entity, error)
	Update(models.Entity, *QEParams) (models.Entity, error)
	Delete(*QEParams) (models.Entity, error)
}

type QEParams struct {
	Selects   []models.IAttribute
	Columns   []models.IAttribute
	Wheres    []models.Condition
	LimitVal  int64
	OffsetVal int64
	OrderBys  []*models.Order
	Populates []models.IAttribute
}

func QEParam() *QEParams {
	return &QEParams{
		[]models.IAttribute{},
		[]models.IAttribute{},
		[]models.Condition{},
		-1,
		0,
		[]*models.Order{},
		[]models.IAttribute{},
	}
}

func (p *QEParams) Select(selects ...models.IAttribute) *QEParams {
	p.Selects = append(p.Selects, selects...)
	return p
}

func (p *QEParams) Insert(columns ...models.IAttribute) *QEParams {
	p.Columns = append(p.Columns, columns...)
	return p
}

func (p *QEParams) Update(columns ...models.IAttribute) *QEParams {
	p.Columns = append(p.Columns, columns...)
	return p
}

func (p *QEParams) Where(where ...models.Condition) *QEParams {
	p.Wheres = append(p.Wheres, where...)
	return p
}

func (p *QEParams) Limit(limit int64) *QEParams {
	p.LimitVal = limit
	return p
}

func (p *QEParams) Offset(offset int64) *QEParams {
	p.OffsetVal = offset
	return p
}

func (p *QEParams) OrderBy(orderBy ...*models.Order) *QEParams {
	p.OrderBys = append(p.OrderBys, orderBy...)
	return p
}

func (p *QEParams) Populate(populate ...models.IAttribute) *QEParams {
	p.Populates = append(p.Populates, populate...)
	return p
}
