package common

type Lifecycle map[string]LifecycleHook

type LifecycleHook func(params DBParams, result any, state any) (afterState any, err error)

func NewLifecycle() Lifecycle {
	return Lifecycle{
		"beforeFindOne":    nullLifecycleHook,
		"beforeFindMany":   nullLifecycleHook,
		"afterFindOne":     nullLifecycleHook,
		"afterFindMany":    nullLifecycleHook,
		"beforeCreate":     nullLifecycleHook,
		"beforeCreateMany": nullLifecycleHook,
		"afterCreate":      nullLifecycleHook,
		"afterCreateMany":  nullLifecycleHook,
		"beforeUpdate":     nullLifecycleHook,
		"beforeUpdateMany": nullLifecycleHook,
		"afterUpdate":      nullLifecycleHook,
		"afterUpdateMany":  nullLifecycleHook,
		"beforeDelete":     nullLifecycleHook,
		"beforeDeleteMany": nullLifecycleHook,
		"afterDelete":      nullLifecycleHook,
		"afterDeleteMany":  nullLifecycleHook,
	}
}

func nullLifecycleHook(params DBParams, result any, state any) (any, error) { return nil, nil }
