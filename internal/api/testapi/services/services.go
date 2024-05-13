package services

import "github.com/cosys-io/cosys/internal/cosys"

var Services = map[string]*cosys.Service{
	"dummy": DummyService,
}
