package common

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Service map[string]any

// A singleton Service map to be used throughout the application
var singleService *Service

var getLock = &sync.Mutex{}
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
func RegisterServices(services Service) (bool, error) {
	serviceSingleton := GetService()
	registerLock.Lock()
	defer registerLock.Unlock()
	for key, value := range services {
		if !isFunc(value) {
			return false, fmt.Errorf("service %s is not a function", key)
		}
		if keyCollision(key, *serviceSingleton) {
			return false, errors.New("service name already exists")
		}
		(*serviceSingleton)[key] = value
	}

	return true, nil
}

// RegisterService registers a service to the singleton Service map
// RegisterService will return an error if the service name already exists in the singleton Service map
// Safe for concurrent use
func RegisterService(serviceName string, service any) (bool, error) {
	if !isFunc(service) {
		return false, fmt.Errorf("service %s is not a function", serviceName)
	}

	serviceSingleton := GetService()
	registerLock.Lock()
	defer registerLock.Unlock()
	if keyCollision(serviceName, *serviceSingleton) {
		return false, errors.New("service name already exists")
	}
	(*serviceSingleton)[serviceName] = service

	return true, nil
}

// GetService returns the singleton Service map
// If the singleton Service map does not exist, it will create one
// Safe for concurrent use
func GetService() *Service {
	if singleService == nil {
		getLock.Lock()
		defer getLock.Unlock()
		if singleService == nil {
			fmt.Println("Creating service instance now.")
			*singleService = make(Service)
		} else {
			fmt.Println("Service instance already created.")
		}
	}
	return singleService
}
