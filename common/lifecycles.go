package common

type LifeCycle map[string]LifeCycleFunc

type LifeCycleFunc func(Event) error

type Event struct {
	Params *DBParams
	Result any
	State  any
}

func NewLifeCycle() LifeCycle {
	return LifeCycle{
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

func noop(event Event) error { return nil }
