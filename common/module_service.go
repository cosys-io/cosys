package common

import (
	"fmt"
	"sync"
)

var (
	msMutex sync.RWMutex
	msMap   = make(map[string]func(*Cosys) ModuleService)
)

func RegisterModuleService(name string, moduleService func(*Cosys) ModuleService) error {
	msMutex.Lock()
	defer msMutex.Unlock()

	if moduleService == nil {
		return fmt.Errorf("module service is nil")
	}

	if _, dup := msMap[name]; dup {
		return fmt.Errorf("duplicate module service:" + name)
	}

	msMap[name] = moduleService
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
