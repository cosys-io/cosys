package common

import (
	"fmt"
	"sync"
)

var (
	dbMutex    sync.RWMutex
	dbRegister = make(map[string]func(*Cosys) Database)
)

func RegisterDatabase(dbName string, dbCtor func(*Cosys) Database) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if dbCtor == nil {
		return fmt.Errorf("database is nil: %s", dbName)
	}

	if _, dup := dbRegister[dbName]; dup {
		return fmt.Errorf("duplicate database: %s", dbName)
	}

	dbRegister[dbName] = dbCtor
	return nil
}

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

type DBParams struct {
	Select   []Attribute
	Columns  []Attribute
	Where    []Condition
	Limit    int64
	Offset   int64
	OrderBy  []*Order
	Populate []Attribute
}

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

type DBParamsBuilder struct {
	selectFields []Attribute
	columns      []Attribute
	where        []Condition
	limit        int64
	offset       int64
	orderBy      []*Order
	populate     []Attribute
}

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

func (p DBParamsBuilder) Select(selects ...Attribute) DBParamsBuilder {
	p.selectFields = append(p.selectFields, selects...)
	return p
}

func (p DBParamsBuilder) Insert(columns ...Attribute) DBParamsBuilder {
	p.columns = append(p.columns, columns...)
	return p
}

func (p DBParamsBuilder) Update(columns ...Attribute) DBParamsBuilder {
	p.columns = append(p.columns, columns...)
	return p
}

func (p DBParamsBuilder) Where(where ...Condition) DBParamsBuilder {
	p.where = append(p.where, where...)
	return p
}

func (p DBParamsBuilder) Limit(limit int64) DBParamsBuilder {
	p.limit = limit
	return p
}

func (p DBParamsBuilder) Offset(offset int64) DBParamsBuilder {
	p.offset = offset
	return p
}

func (p DBParamsBuilder) OrderBy(orderBy ...*Order) DBParamsBuilder {
	p.orderBy = append(p.orderBy, orderBy...)
	return p
}

func (p DBParamsBuilder) Populate(populate ...Attribute) DBParamsBuilder {
	p.populate = append(p.populate, populate...)
	return p
}

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
