package common

import (
	"fmt"
	"maps"
	"reflect"
	"sync"
)

// Config

// options are configurations for registers.
type options struct {
	itemName    string
	checkZero   bool
	allowUpdate bool
	allowRemove bool
	mustExist   bool
}

// defaultOptions returns the default configuration.
func defaultOptions() options {
	return options{
		itemName:    "item",
		checkZero:   true,
		allowUpdate: false,
		allowRemove: false,
		mustExist:   false,
	}
}

// option is a configuration.
type option func(*options)

// itemName sets the item name used in error messages.
func itemName(itemName string) option {
	return func(opts *options) {
		opts.itemName = itemName
	}
}

// skipCheckZero skips zero/nil checks.
// skipCheckZero would typically be used for non-pointer types,
// where the zero value could be valid.
var skipCheckZero option = func(opts *options) {
	opts.checkZero = false
}

// allowUpdate allows the register item to be updated.
var allowUpdate option = func(opts *options) {
	opts.allowUpdate = true
}

// allowRemove allows the register item to be removed.
var allowRemove option = func(opts *options) {
	opts.allowRemove = true
}

// mustExist will cause an error to be thrown when updating or removing an unregistered value.
var mustExist option = func(opts *options) {
	opts.mustExist = true
}

// Single Register

// singleRegister is a register for a single value.
type singleRegister[T any] struct {
	mutex      sync.RWMutex
	register   T
	registered bool
	options    options
}

// newSingleRegister returns a new single register with configurations.
func newSingleRegister[T any](cfg ...option) *singleRegister[T] {
	opts := defaultOptions()
	for _, opt := range cfg {
		opt(&opts)
	}

	return &singleRegister[T]{
		mutex:      sync.RWMutex{},
		register:   *new(T),
		registered: false,
		options:    opts,
	}
}

// Get returns the registered value and throws an error
// if a value has not been registered.
// Safe for concurrent use.
func (r *singleRegister[T]) Get() (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if !r.registered {
		return *new(T), fmt.Errorf("%s not registered", r.options.itemName)
	}

	return r.register, nil
}

// Clone returns a cloned single register
// Safe for concurrent use.
func (r *singleRegister[T]) Clone() *singleRegister[T] {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return &singleRegister[T]{
		mutex:      sync.RWMutex{},
		register:   r.register,
		registered: r.registered,

		options: r.options,
	}
}

// Register sets the registered value, and throws an error
// if the checkZero configuration is true and the value is a zero-value,
// or if a value has already been registered.
// Safe for concurrent use.
func (r *singleRegister[T]) Register(item T) error {
	if r.options.checkZero && reflect.DeepEqual(item, *new(T)) {
		return fmt.Errorf("%s is nil", r.options.itemName)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.registered {
		return fmt.Errorf("%s is already registered", r.options.itemName)
	}

	r.register = item
	r.registered = true

	return nil
}

// Update sets a new registered value, and throws an error
// if allowUpdate configuration is false,
// or if the checkZero configuration is true and the value is a zero-value,
// or if the mustExist configuration is true and a value has not been registered or has been removed.
// Safe for concurrent use.
func (r *singleRegister[T]) Update(item T) error {
	if !r.options.allowUpdate {
		return fmt.Errorf("%s cannot be updated", r.options.itemName)
	}

	if r.options.checkZero && reflect.DeepEqual(item, *new(T)) {
		return fmt.Errorf("%s is nil", r.options.itemName)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.options.mustExist && !r.registered {
		return fmt.Errorf("%s is not registered", r.options.itemName)
	}

	r.register = item

	return nil
}

// Remove deletes the registered value, and throws an error
// if the allowRemove configuration is false,
// or if the mustExists configuration is true and a value has not been registered or has been removed.
// Safe for concurrent use.
func (r *singleRegister[T]) Remove() error {
	if !r.options.allowRemove {
		return fmt.Errorf("%s cannot be removed", r.options.itemName)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.options.mustExist && !r.registered {
		return fmt.Errorf("%s is not registered", r.options.itemName)
	}

	r.register = *new(T)
	r.registered = false

	return nil
}

// Multi Register

// multiRegister is a register for multiple values.
type multiRegister[T any] struct {
	mutex    sync.RWMutex
	register map[string]T
	options  options
}

// newMultiRegister returns a new multi register with configurations.
func newMultiRegister[T any](cfg ...option) *multiRegister[T] {
	opts := defaultOptions()
	for _, opt := range cfg {
		opt(&opts)
	}

	return &multiRegister[T]{
		mutex:    sync.RWMutex{},
		register: make(map[string]T),
		options:  opts,
	}
}

// Get return the registered value with the given uid, and throws an error
// if a value has not been registered under that uid.
func (r *multiRegister[T]) Get(uid string) (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	v, ok := r.register[uid]
	if !ok {
		return *new(T), fmt.Errorf("%s not found: %s", r.options.itemName, uid)
	}

	return v, nil
}

// GetAll returns a map of all registered values.
// Safe for concurrent use, returns a copy of underlying map.
func (r *multiRegister[T]) GetAll() map[string]T {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return maps.Clone(r.register)
}

// Clone returns a cloned multi register.
// Safe for concurrent use.
func (r *multiRegister[T]) Clone() *multiRegister[T] {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return &multiRegister[T]{
		mutex:    sync.RWMutex{},
		register: maps.Clone(r.register),

		options: r.options,
	}
}

// Register sets a value under the given uid, and throws an error
// if the checkZero configuration is true and the value is a zero-value,
// of if a value has been registered under that uid.
// Safe for concurrent use.
func (r *multiRegister[T]) Register(uid string, item T) error {
	if r.options.checkZero && reflect.DeepEqual(item, *new(T)) {
		return fmt.Errorf("%s is nil: %s", r.options.itemName, uid)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, dup := r.register[uid]; dup {
		return fmt.Errorf("duplicate %s: %s", r.options.itemName, uid)
	}

	r.register[uid] = item

	return nil
}

// Update sets a new value under the given uid, and throws an error
// if the allowUpdate configuration is false,
// or if the checkZero configuration is true and the value is a zero-value,
// or if the mustExist configuration is true and no value has been registered under that uid.
// Safe for concurrent use.
func (r *multiRegister[T]) Update(uid string, item T) error {
	if !r.options.allowUpdate {
		return fmt.Errorf("%s cannot be updated", r.options.itemName)
	}

	if r.options.checkZero && reflect.DeepEqual(item, *new(T)) {
		return fmt.Errorf("%s is nil: %s", r.options.itemName, uid)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.register[uid]; r.options.mustExist && !ok {
		return fmt.Errorf("%s not found: %s", r.options.itemName, uid)
	}

	r.register[uid] = item

	return nil
}

// Remove deletes the value under the given uid, and throws an error
// if the allowRemove configuration is false,
// or if the mustExist configuration is true and no value has been registered under that uid.
// Safe for concurrent use.
func (r *multiRegister[T]) Remove(uid string) error {
	if !r.options.allowRemove {
		return fmt.Errorf("%s cannot be removed", r.options.itemName)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.register[uid]; r.options.mustExist && !ok {
		return fmt.Errorf("%s not exist: %s", r.options.itemName, uid)
	}

	delete(r.register, uid)

	return nil
}
