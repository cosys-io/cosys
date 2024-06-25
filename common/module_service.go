package common

import (
	"fmt"
	"sync"
)

var (
	msMutex    sync.RWMutex
	msRegister = make(map[string]func(*Cosys) ModuleService)
)

func RegisterModuleService(msName string, msCtor func(*Cosys) ModuleService) error {
	msMutex.Lock()
	defer msMutex.Unlock()

	if msCtor == nil {
		return fmt.Errorf("module service is nil: %s", msName)
	}

	if _, dup := msRegister[msName]; dup {
		return fmt.Errorf("duplicate module service: %s", msName)
	}

	msRegister[msName] = msCtor
	return nil
}

type ModuleService interface {
	FindOne(uid string, id int, params MSParams) (Entity, error)
	FindMany(uid string, params MSParams) ([]Entity, error)
	Create(uid string, data Entity, params MSParams) (Entity, error)
	Update(uid string, data Entity, id int, params MSParams) (Entity, error)
	Delete(uid string, id int, params MSParams) (Entity, error)
}

type MSParams struct {
	Select   []Attribute
	Fields   []Attribute
	Filter   []Condition
	Start    int64
	Limit    int64
	Sort     []*Order
	Populate []Attribute
}

func NewMSParams() MSParams {
	return MSParams{
		Select:   []Attribute{},
		Fields:   []Attribute{},
		Filter:   []Condition{},
		Start:    0,
		Limit:    -1,
		Sort:     []*Order{},
		Populate: []Attribute{},
	}
}

type MSParamsBuilder struct {
	selectFields []Attribute
	fields       []Attribute
	filter       []Condition
	start        int64
	limit        int64
	sort         []*Order
	populate     []Attribute
}

func NewMSParamsBuilder() MSParamsBuilder {
	return MSParamsBuilder{
		[]Attribute{},
		[]Attribute{},
		[]Condition{},
		0,
		-1,
		[]*Order{},
		[]Attribute{},
	}
}

func (p MSParamsBuilder) GetField(fields ...Attribute) MSParamsBuilder {
	p.selectFields = append(p.selectFields, fields...)
	return p
}

func (p MSParamsBuilder) SetField(fields ...Attribute) MSParamsBuilder {
	p.fields = append(p.fields, fields...)
	return p
}

func (p MSParamsBuilder) Filter(filters ...Condition) MSParamsBuilder {
	p.filter = append(p.filter, filters...)
	return p
}

func (p MSParamsBuilder) Limit(limit int64) MSParamsBuilder {
	p.limit = limit
	return p
}

func (p MSParamsBuilder) Start(start int64) MSParamsBuilder {
	p.start = start
	return p
}

func (p MSParamsBuilder) Sort(sorts ...*Order) MSParamsBuilder {
	p.sort = append(p.sort, sorts...)
	return p
}

func (p MSParamsBuilder) Populate(populates ...Attribute) MSParamsBuilder {
	p.populate = append(p.populate, populates...)
	return p
}

func (p MSParamsBuilder) Build() MSParams {
	return MSParams{
		Select:   p.selectFields,
		Fields:   p.fields,
		Filter:   p.filter,
		Start:    p.start,
		Limit:    p.limit,
		Sort:     p.sort,
		Populate: p.populate,
	}
}
