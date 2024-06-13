package common

type Lifecycle map[string]LifecycleFunc

type LifecycleFunc func(*Event) error

type Event struct {
	Params DBParams
	Result any
	State  any
}

func NewLifeCycle() Lifecycle {
	return Lifecycle{
		"beforeFindOne":    noop,
		"beforeFindMany":   noop,
		"afterFindOne":     noop,
		"afterFindMany":    noop,
		"beforeCreate":     noop,
		"beforeCreateMany": noop,
		"afterCreate":      noop,
		"afterCreateMany":  noop,
		"beforeUpdate":     noop,
		"beforeUpdateMany": noop,
		"afterUpdate":      noop,
		"afterUpdateMany":  noop,
		"beforeDelete":     noop,
		"beforeDeleteMany": noop,
		"afterDelete":      noop,
		"afterDeleteMany":  noop,
	}
}

func NewEvent(params DBParams) Event {
	return Event{
		Params: params,
		Result: nil,
		State:  nil,
	}
}

func noop(event *Event) error { return nil }
