package common

import (
	"fmt"
	"maps"
	"math/rand"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

// Config

// options are configurations for registers.
type options struct {
	itemName   string
	checkZero  bool
	checkExist bool
}

// defaultOptions returns the default configuration.
func defaultOptions() options {
	return options{
		itemName:   "item",
		checkZero:  true,
		checkExist: true,
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

// skipCheckExist will cause an error to be thrown when updating or removing an unregistered value.
var skipCheckExist option = func(opts *options) {
	opts.checkExist = true
}

// Single Register

// singleRegister is a register for a single value.
type singleRegister[T any] struct {
	mutex      *sync.RWMutex
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
		mutex:      &sync.RWMutex{},
		register:   *new(T),
		registered: false,
		options:    opts,
	}
}

// Get returns the registered value and returns an error
// if a value has not been registered.
// Safe for concurrent use.
func (r singleRegister[T]) Get() (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if !r.registered {
		return *new(T), fmt.Errorf("%s not registered", r.options.itemName)
	}

	return r.register, nil
}

// Clone returns a cloned single register
// Safe for concurrent use.
func (r singleRegister[T]) Clone() *singleRegister[T] {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return &singleRegister[T]{
		mutex:      &sync.RWMutex{},
		register:   r.register,
		registered: r.registered,

		options: r.options,
	}
}

// Register sets the registered value, and returns an error
// if the checkZero configuration is true and the value is a zero-value,
// or if a value has already been registered.
// Safe for concurrent use.
func (r singleRegister[T]) Register(item T) error {
	if r.options.checkZero && isZero(item) {
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

// Permanent Register

// permRegister is a register for multiple values which cannot be updated or deleted.
type permRegister[T any] struct {
	mutex    *sync.RWMutex
	register map[string]T
	options  options
}

// newPermRegister returns a new permanent register with configurations.
func newPermRegister[T any](cfg ...option) *permRegister[T] {
	opts := defaultOptions()
	for _, opt := range cfg {
		opt(&opts)
	}

	return &permRegister[T]{
		mutex:    &sync.RWMutex{},
		register: make(map[string]T),
		options:  opts,
	}
}

// Get return the registered value with the given uid, and returns an error
// if a value has not been registered under that uid.
func (r permRegister[T]) Get(uid string) (T, error) {
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
func (r permRegister[T]) GetAll() map[string]T {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return maps.Clone(r.register)
}

// Clone returns a cloned multi register.
// Safe for concurrent use.
func (r permRegister[T]) Clone() *permRegister[T] {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return &permRegister[T]{
		mutex:    &sync.RWMutex{},
		register: maps.Clone(r.register),

		options: r.options,
	}
}

// Register sets a value under the given uid, and returns an error
// if the checkZero configuration is true and the value is a zero-value,
// of if a value has been registered under that uid.
// Safe for concurrent use.
func (r permRegister[T]) Register(uid string, item T) error {
	if r.options.checkZero && isZero(item) {
		return fmt.Errorf("%s is nil: %s", r.options.itemName, uid)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if isDup(r.register, uid) {
		return fmt.Errorf("duplicate %s: %s", r.options.itemName, uid)
	}

	r.register[uid] = item

	return nil
}

// RegisterMany sets a map of values under their corresponding uid, and returns an error
// if the checkZero configuration is true and any value is a zero-value
// or if a value has been registered under any uid.
// The operation is atomic, either all or no values will be set.
// Safe for concurrent use.
func (r permRegister[T]) RegisterMany(items map[string]T) error {
	if r.options.checkZero {
		if err := anyZero(r.options.itemName, items); err != nil {
			return err
		}
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err := anyDup(r.options.itemName, r.register, items); err != nil {
		return err
	}

	for uid, item := range items {
		r.register[uid] = item
	}

	return nil
}

// RegisterRandom sets a value under a random uid, and returns the uid or returns an error
// if the checkZero configuration is true and the value is a zero-value,
// or if it fails to generate a valid random uid.
// The generated uid are all prefixed with '$', if using both Register and RegisterRandom,
// do not use uid prefixed by $ for Register.
// Safe for concurrent use.
func (r permRegister[T]) RegisterRandom(item T) (string, error) {
	if r.options.checkZero && isZero(item) {
		return "", fmt.Errorf("%s is nil", r.options.itemName)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	var uid string
	for _ = range 1000 {
		uid = randomString([]byte{'$'}, 8)
		if isDup(r.register, uid) {
			continue
		}
	}

	if uid == "" {
		return "", fmt.Errorf("could not generate uid for %s", r.options.itemName)
	}

	r.register[uid] = item

	return uid, nil
}

// Multi Register

// multiRegister is a register for multiple values.
type multiRegister[T any] struct {
	*permRegister[T]
}

// newMultiRegister returns a new multi register with configurations.
func newMultiRegister[T any](cfg ...option) *multiRegister[T] {
	return &multiRegister[T]{
		newPermRegister[T](cfg...),
	}
}

// Clone returns a cloned multi register.
// Safe for concurrent use.
func (r multiRegister[T]) Clone() *multiRegister[T] {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return &multiRegister[T]{
		r.permRegister.Clone(),
	}
}

// Update sets a new value under the given uid, and returns an error
// if the allowUpdate configuration is false,
// or if the checkZero configuration is true and the value is a zero-value,
// or if the mustExist configuration is true and no value has been registered under that uid.
// Safe for concurrent use.
func (r multiRegister[T]) Update(uid string, item T) error {
	if r.options.checkZero && isZero(item) {
		return fmt.Errorf("%s is nil: %s", r.options.itemName, uid)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.options.checkExist && isMissing(r.register, uid) {
		return fmt.Errorf("%s not found: %s", r.options.itemName, uid)
	}

	r.register[uid] = item

	return nil
}

// Remove deletes the value under the given uid, and returns an error
// if the allowRemove configuration is false,
// or if the mustExist configuration is true and no value has been registered under that uid.
// Safe for concurrent use.
func (r multiRegister[T]) Remove(uid string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.options.checkExist && isMissing(r.register, uid) {
		return fmt.Errorf("%s not exist: %s", r.options.itemName, uid)
	}

	delete(r.register, uid)

	return nil
}

// Stringer Register

// stringerRegister is a register of stringer values.
type stringerRegister[T fmt.Stringer] struct {
	*multiRegister[T]
}

// newStringerRegister returns a new stringer register with configurations.
func newStringerRegister[T fmt.Stringer](cfg ...option) *stringerRegister[T] {
	return &stringerRegister[T]{
		newMultiRegister[T](cfg...),
	}
}

// RegisterStringers sets a slice of stringer values under their respective String() return values,
// and returns an error if any values have the same String() return values,
// or if the checkZero configuration is true and any value is a zero-value,
// or if a value has been registered under any String() return value.
// The operation is atomic, either all or no values will be set.
// Safe for concurrent use.
func (r stringerRegister[T]) RegisterStringers(items ...T) error {
	itemMap, err := toMap(r.options.itemName, items)
	if err != nil {
		return err
	}

	if r.options.checkZero {
		if err = anyZero(r.options.itemName, itemMap); err != nil {
			return err
		}
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if err = anyDup(r.options.itemName, r.register, itemMap); err != nil {
		return err
	}

	for uid, item := range itemMap {
		r.register[uid] = item
	}

	return nil
}

func (r stringerRegister[T]) Clone() *stringerRegister[T] {
	return &stringerRegister[T]{
		r.multiRegister.Clone(),
	}
}

// Utils

// chars are the possible chars used in randomString
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	charsIndexLen  = 6                    // 6 bits to represent an index in chars
	charsIndexMask = 1<<charsIndexLen - 1 // 0b111111 bit mask
	charsIndexMax  = 63 / charsIndexLen   // # of char indices fitting in 63 bits
)

// src is a random generator
var src = rand.NewSource(time.Now().UnixNano())

// randomString generates a random string of length n
func randomString(prefix []byte, num int) string {
	prefixLen := len(prefix)
	buffer := make([]byte, prefixLen+num)
	copy(buffer, prefix)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := prefixLen+num-1, src.Int63(), charsIndexMax; i >= prefixLen; {
		if remain == 0 {
			cache, remain = src.Int63(), charsIndexMax
		}
		if idx := int(cache & charsIndexMask); idx < len(chars) {
			buffer[i] = chars[idx]
			i--
		}
		cache >>= charsIndexLen
		remain--
	}

	return *(*string)(unsafe.Pointer(&buffer))
}

// isZero returns whether an item is a zero-value.
func isZero[T any](item T) bool {
	return reflect.ValueOf(item).IsZero()
}

// anyZero returns whether any item in a map is a zero-value.
func anyZero[T any](itemName string, items map[string]T) error {
	for uid, item := range items {
		if isZero(item) {
			return fmt.Errorf("%s is nil: %s", itemName, uid)
		}
	}

	return nil
}

// isDup returns whether a map already contains a key.
func isDup[T any](items map[string]T, uid string) bool {
	_, dup := items[uid]
	return dup
}

// anyDup returns whether any keys in two maps are the same.
func anyDup[T any](itemName string, items1 map[string]T, items2 map[string]T) error {
	if len(items2) > len(items1) {
		items1, items2 = items2, items1
	}

	for uid := range items2 {
		if isDup(items1, uid) {
			return fmt.Errorf("duplicate %s: %s", itemName, uid)
		}
	}

	return nil
}

// isMissing returns if a key is missing from a map.
func isMissing[T any](items map[string]T, uid string) bool {
	_, exists := items[uid]
	return !exists
}

// isStringer returns whether a type parameter is a Stringer.
func isStringer[T any]() bool {
	_, ok := reflect.Zero(reflect.TypeOf(*new(T))).Interface().(fmt.Stringer)
	return ok
}

// toMap converts a generic slice into a map, and returns an error
// if any entry is not a Stringer.
func toMap[T any](itemName string, itemSlice []T) (map[string]T, error) {
	itemMap := map[string]T{}
	for _, item := range itemSlice {
		stringer, ok := reflect.ValueOf(item).Interface().(fmt.Stringer)
		if !ok {
			return nil, fmt.Errorf("%s is not stringer", itemName)
		}
		uid := stringer.String()

		if isDup(itemMap, uid) {
			return nil, fmt.Errorf("duplicate %s: %s", itemName, uid)
		}

		itemMap[uid] = item
	}

	return itemMap, nil
}
