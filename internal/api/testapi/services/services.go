package services

import "github.com/cosys-io/cosys/internal/common"

var Services = map[string]*common.Service{
	"dummy": DummyService,
}
