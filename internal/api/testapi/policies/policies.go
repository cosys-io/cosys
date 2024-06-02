package policies

import "github.com/cosys-io/cosys/internal/apis"

var Policies = map[string]apis.Policy{
	"dummy": DummyPolicy,
}
