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
	Selects   []Attribute
	Columns   []Attribute
	Wheres    []Condition
	LimitVal  int64
	OffsetVal int64
	OrderBys  []*Order
	Populates []Attribute
}

func DBParam() DBParams {
	return DBParams{
		[]Attribute{},
		[]Attribute{},
		[]Condition{},
		-1,
		0,
		[]*Order{},
		[]Attribute{},
	}
}

func (p DBParams) Select(selects ...Attribute) DBParams {
	p.Selects = append(p.Selects, selects...)
	return p
}

func (p DBParams) Insert(columns ...Attribute) DBParams {
	p.Columns = append(p.Columns, columns...)
	return p
}

func (p DBParams) Update(columns ...Attribute) DBParams {
	p.Columns = append(p.Columns, columns...)
	return p
}

func (p DBParams) Where(where ...Condition) DBParams {
	p.Wheres = append(p.Wheres, where...)
	return p
}

func (p DBParams) Limit(limit int64) DBParams {
	p.LimitVal = limit
	return p
}

func (p DBParams) Offset(offset int64) DBParams {
	p.OffsetVal = offset
	return p
}

func (p DBParams) OrderBy(orderBy ...*Order) DBParams {
	p.OrderBys = append(p.OrderBys, orderBy...)
	return p
}

func (p DBParams) Populate(populate ...Attribute) DBParams {
	p.Populates = append(p.Populates, populate...)
	return p
}
