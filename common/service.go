package common

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Service map[string]any

// A singleton Service map to be used throughout the application
var singleService Service = make(Service)

var registerLock = &sync.Mutex{}

func isFunc(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Func
}

func keyCollision(key string, service Service) bool {
	_, ok := service[key]
	return ok
}

// RegisterServices registers a map of services to the singleton Service map
// RegisterServices will return an error if a service name already exists in the singleton Service map
// Safe for concurrent use
func RegisterServices(services Service) error {
	serviceSingleton := GetService()
	registerLock.Lock()
	defer registerLock.Unlock()
	for key, value := range services {
		if !isFunc(value) {
			return fmt.Errorf("service %s is not a function", key)
		}
		if keyCollision(key, serviceSingleton) {
			return errors.New("service name already exists")
		}
		serviceSingleton[key] = value
	}

	return nil
}

// RegisterService registers a service to the singleton Service map
// RegisterService will return an error if the service name already exists in the singleton Service map
// Safe for concurrent use
func RegisterService(serviceName string, service any) error {
	if !isFunc(service) {
		return fmt.Errorf("service %s is not a function", serviceName)
	}

	serviceSingleton := GetService()
	registerLock.Lock()
	defer registerLock.Unlock()
	if keyCollision(serviceName, serviceSingleton) {
		return errors.New("service name already exists")
	}
	serviceSingleton[serviceName] = service

	return nil
}

// GetService returns the singleton Service map
// Unsafe for concurrent use as it exposes the underlying singleService object to everyone
func GetService() Service {
	return singleService
}
