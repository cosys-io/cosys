package common

type Lifecycle map[string]LifecycleFunc

type LifecycleFunc func(params DBParams, result any, state any) (afterState any, err error)

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

func noop(params DBParams, result any, state any) (any, error) { return nil, nil }
