package common

import (
	"fmt"
	"maps"
	"reflect"
	"sync"
)

// Single Register

type singleRegister[T any] struct {
	mutex      sync.RWMutex
	register   *T
	registered bool

	itemName    string
	checkZero   bool
	allowUpdate bool
}

func newSingleRegister[T any](itemName string, checkZero bool, allowUpdate bool) *singleRegister[T] {
	return &singleRegister[T]{
		mutex:      sync.RWMutex{},
		register:   nil,
		registered: false,

		itemName:    itemName,
		checkZero:   checkZero,
		allowUpdate: allowUpdate,
	}
}

func (r *singleRegister[T]) Get() (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if !r.registered {
		return *new(T), fmt.Errorf("%s not registered", r.itemName)
	}

	return *r.register, nil
}

func (r *singleRegister[T]) Clone() *singleRegister[T] {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return &singleRegister[T]{
		mutex:      sync.RWMutex{},
		register:   r.register,
		registered: r.registered,

		itemName:    r.itemName,
		checkZero:   r.checkZero,
		allowUpdate: r.allowUpdate,
	}
}

func (r *singleRegister[T]) Register(item *T) error {
	if item == nil {
		return fmt.Errorf("%s is nil", r.itemName)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.allowUpdate && r.registered {
		return fmt.Errorf("%s is already registered", r.itemName)
	}

	r.register = item
	r.registered = true

	return nil
}

// Multi Register

type multiRegister[T any] struct {
	mutex    sync.RWMutex
	register map[string]T

	itemName  string
	checkZero bool
}

func newMultiRegister[T any](itemName string, checkZero bool) *multiRegister[T] {
	return &multiRegister[T]{
		mutex:    sync.RWMutex{},
		register: make(map[string]T),

		itemName:  itemName,
		checkZero: checkZero,
	}
}

func (r *multiRegister[T]) Get(uid string) (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	v, ok := r.register[uid]
	if !ok {
		return *new(T), fmt.Errorf("%s not found: %s", r.itemName, uid)
	}

	return v, nil
}

func (r *multiRegister[T]) GetAll() map[string]T {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return maps.Clone(r.register)
}

func (r *multiRegister[T]) Clone() *multiRegister[T] {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return &multiRegister[T]{
		mutex:    sync.RWMutex{},
		register: maps.Clone(r.register),

		itemName:  r.itemName,
		checkZero: r.checkZero,
	}
}

func (r *multiRegister[T]) Register(uid string, item T) error {
	if r.checkZero && reflect.DeepEqual(item, reflect.Zero(reflect.TypeOf(item))) {
		return fmt.Errorf("%s is nil: %s", r.itemName, uid)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, dup := r.register[uid]; dup {
		return fmt.Errorf("duplicate %s: %s", r.itemName, uid)
	}

	r.register[uid] = item

	return nil
}

func (r *multiRegister[T]) Update(uid string, item T) error {
	if r.checkZero && reflect.DeepEqual(item, reflect.Zero(reflect.TypeOf(item))) {
		return fmt.Errorf("%s is nil: %s", r.itemName, uid)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.register[uid]; !ok {
		return fmt.Errorf("%s not found: %s", r.itemName, uid)
	}

	r.register[uid] = item

	return nil
}

func (r *multiRegister[T]) Remove(uid string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.register[uid]; !ok {
		return fmt.Errorf("%s not exist: %s", r.itemName, uid)
	}

	delete(r.register, uid)

	return nil
}
