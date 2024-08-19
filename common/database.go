package common

// Database is a core service for interacting with the relational database.
type Database interface {
	FindOne(uid string, params DBParams) (Entity, error)
	FindMany(uid string, params DBParams) ([]Entity, error)
	Create(uid string, data Entity, params DBParams) (Entity, error)
	CreateMany(uid string, data []Entity, params DBParams) ([]Entity, error)
	Update(uid string, data Entity, params DBParams) (Entity, error)
	UpdateMany(uid string, data Entity, params DBParams) ([]Entity, error)
	Delete(uid string, params DBParams) (Entity, error)
	DeleteMany(uid string, params DBParams) ([]Entity, error)
}

// DBParams are query conditions.
type DBParams struct {
	Select   []Attribute
	Columns  []Attribute
	Where    []Condition
	Limit    int64
	Offset   int64
	OrderBy  []*Order
	Populate []Attribute
}

// NewDBParams returns a new DBParams with default conditions.
func NewDBParams() DBParams {
	return DBParams{
		Select:   []Attribute{},
		Columns:  []Attribute{},
		Where:    []Condition{},
		Limit:    -1,
		Offset:   0,
		OrderBy:  []*Order{},
		Populate: []Attribute{},
	}
}

// DBParamsBuilder is a builder for DBParams.
type DBParamsBuilder struct {
	selectFields []Attribute
	columns      []Attribute
	where        []Condition
	limit        int64
	offset       int64
	orderBy      []*Order
	populate     []Attribute
}

// NewDBParamsBuilder returns a new DBParamsBuilder with default conditions.
func NewDBParamsBuilder() DBParamsBuilder {
	return DBParamsBuilder{
		[]Attribute{},
		[]Attribute{},
		[]Condition{},
		-1,
		0,
		[]*Order{},
		[]Attribute{},
	}
}

// Select adds attributes to return for all queries.
func (p DBParamsBuilder) Select(selects ...Attribute) DBParamsBuilder {
	p.selectFields = append(p.selectFields, selects...)
	return p
}

// Insert adds attributes to insert for insert queries.
func (p DBParamsBuilder) Insert(columns ...Attribute) DBParamsBuilder {
	p.columns = append(p.columns, columns...)
	return p
}

// Update adds attributes to update for update queries.
func (p DBParamsBuilder) Update(columns ...Attribute) DBParamsBuilder {
	p.columns = append(p.columns, columns...)
	return p
}

// Where adds where conditions for all queries.
func (p DBParamsBuilder) Where(where ...Condition) DBParamsBuilder {
	p.where = append(p.where, where...)
	return p
}

// Limit sets the limit condition for all queries.
func (p DBParamsBuilder) Limit(limit int64) DBParamsBuilder {
	p.limit = limit
	return p
}

// Offset sets the offset condition for all queries.
func (p DBParamsBuilder) Offset(offset int64) DBParamsBuilder {
	p.offset = offset
	return p
}

// OrderBy adds order-by conditions for all queries.
func (p DBParamsBuilder) OrderBy(orderBy ...*Order) DBParamsBuilder {
	p.orderBy = append(p.orderBy, orderBy...)
	return p
}

// Populate adds populate conditions for all queries.
func (p DBParamsBuilder) Populate(populate ...Attribute) DBParamsBuilder {
	p.populate = append(p.populate, populate...)
	return p
}

// Build returns the DBParams with the set conditions.
func (p DBParamsBuilder) Build() DBParams {
	return DBParams{
		Select:   p.selectFields,
		Columns:  p.columns,
		Where:    p.where,
		Limit:    p.limit,
		Offset:   p.offset,
		OrderBy:  p.orderBy,
		Populate: p.populate,
	}
}
