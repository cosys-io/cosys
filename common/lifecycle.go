package common

import "fmt"

// LifecycleHook is a hook that is called when a model's lifecycle event happens.
type LifecycleHook func(params DBParams, result any, state any) (afterState any, err error)

// noopLifecycleHook is a hook that does nothing.
func noopLifecycleHook(params DBParams, result any, state any) (any, error) { return nil, nil }

// Lifecycle is a group of lifecycle hooks associated with a model.
type Lifecycle struct {
	beforeFindOne    *multiRegister[LifecycleHook]
	afterFindOne     *multiRegister[LifecycleHook]
	beforeFindMany   *multiRegister[LifecycleHook]
	afterFindMany    *multiRegister[LifecycleHook]
	beforeCreate     *multiRegister[LifecycleHook]
	afterCreate      *multiRegister[LifecycleHook]
	beforeCreateMany *multiRegister[LifecycleHook]
	afterCreateMany  *multiRegister[LifecycleHook]
	beforeUpdate     *multiRegister[LifecycleHook]
	afterUpdate      *multiRegister[LifecycleHook]
	beforeUpdateMany *multiRegister[LifecycleHook]
	afterUpdateMany  *multiRegister[LifecycleHook]
	beforeDelete     *multiRegister[LifecycleHook]
	afterDelete      *multiRegister[LifecycleHook]
	beforeDeleteMany *multiRegister[LifecycleHook]
	afterDeleteMany  *multiRegister[LifecycleHook]
}

// getRegister returns the register corresponding to the lifecycle event.
func (l Lifecycle) getRegister(event string) (*multiRegister[LifecycleHook], error) {
	switch event {
	case "beforeFindOne":
		return l.beforeFindOne, nil
	case "afterFindOne":
		return l.afterFindOne, nil
	case "beforeFindMany":
		return l.beforeFindMany, nil
	case "afterFindMany":
		return l.afterFindMany, nil
	case "beforeCreate":
		return l.beforeCreate, nil
	case "afterCreate":
		return l.afterCreate, nil
	case "beforeCreateMany":
		return l.beforeCreateMany, nil
	case "afterCreateMany":
		return l.afterCreateMany, nil
	case "beforeUpdate":
		return l.beforeUpdate, nil
	case "afterUpdate":
		return l.afterUpdate, nil
	case "beforeUpdateMany":
		return l.beforeUpdateMany, nil
	case "afterUpdateMany":
		return l.afterUpdateMany, nil
	case "beforeDelete":
		return l.beforeDelete, nil
	case "afterDelete":
		return l.afterDelete, nil
	case "beforeDeleteMany":
		return l.beforeDeleteMany, nil
	case "afterDeleteMany":
		return l.afterDeleteMany, nil
	default:
		return nil, fmt.Errorf("unknown lifecycle event: %s", event)
	}
}

// Add adds a hook for a lifecycle event and returns a uid used for updating and removing.
// Safe for concurrent use.
func (l Lifecycle) Add(event string, hook LifecycleHook) (string, error) {
	register, err := l.getRegister(event)
	if err != nil {
		return "", err
	}

	return register.RegisterRandom(hook)
}

// Update updates a hook specified by its uid for a lifecycle event.
// Safe for concurrent use.
func (l Lifecycle) Update(event string, uid string, hook LifecycleHook) error {
	register, err := l.getRegister(event)
	if err != nil {
		return err
	}

	return register.Update(uid, hook)
}

// Remove removes a hook specified by its uid for a lifecycle event.
// Safe for concurrent use.
func (l Lifecycle) Remove(event string, uid string) error {
	register, err := l.getRegister(event)
	if err != nil {
		return err
	}

	return register.Remove(uid)
}

// NewLifecycle returns a new lifecycle.
func NewLifecycle() Lifecycle {
	return Lifecycle{
		beforeFindOne:    newMultiRegister[LifecycleHook](),
		afterFindOne:     newMultiRegister[LifecycleHook](),
		beforeFindMany:   newMultiRegister[LifecycleHook](),
		afterFindMany:    newMultiRegister[LifecycleHook](),
		beforeCreate:     newMultiRegister[LifecycleHook](),
		afterCreate:      newMultiRegister[LifecycleHook](),
		beforeCreateMany: newMultiRegister[LifecycleHook](),
		afterCreateMany:  newMultiRegister[LifecycleHook](),
		beforeUpdate:     newMultiRegister[LifecycleHook](),
		afterUpdate:      newMultiRegister[LifecycleHook](),
		beforeUpdateMany: newMultiRegister[LifecycleHook](),
		afterUpdateMany:  newMultiRegister[LifecycleHook](),
		beforeDelete:     newMultiRegister[LifecycleHook](),
		afterDelete:      newMultiRegister[LifecycleHook](),
		beforeDeleteMany: newMultiRegister[LifecycleHook](),
		afterDeleteMany:  newMultiRegister[LifecycleHook](),
	}
}
