package common

type Module struct {
	Routes      []*Route
	Controllers map[string]Controller
	Middlewares map[string]Middleware
	Policies    map[string]Policy

	Models   map[string]Model
	Services map[string]Service

	OnRegister func(Cosys) (Cosys, error)
	OnDestroy  func(Cosys) (Cosys, error)
}
