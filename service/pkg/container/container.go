package container

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Represents the dependency injection container.
// It allows to register and resolve dependencies via the Register and Get functions.
// Thread safety access provided via sync.Mutex.
type Container struct {
	registry map[string]interface{}
	mutex    sync.Mutex
}

// Creates a new Container.
// The Container represents the dependency injection container.
// It allows to register and resolve dependencies via the Register and Get functions.
// Thread safety access provided via sync.Mutex.
func NewContainer() *Container {
	return &Container{
		registry: make(map[string]interface{}),
	}
}

func (c *Container) register(key string, dependency interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.registry[key]; exists {
		return errors.New(fmt.Sprintf("dependency %s already registered", key))
	}
	dependencyType := reflect.TypeOf(dependency)
	if dependencyType.Kind() != reflect.Func || dependencyType.NumIn() != 0 || dependencyType.NumOut() != 1 {
		return errors.New(fmt.Sprintf("dependency factory must take no arguments and return a single output"))
	}
	c.registry[key] = dependency
	return nil
}

func (c *Container) get(ptr interface{}, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	dependency, exists := c.registry[key]
	if !exists {
		return errors.New(fmt.Sprintf("dependency %s not found in container", key))
	}

	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("expected non nil pointer")
	}

	dependencyVal := reflect.ValueOf(dependency)
	resolved := dependencyVal.Call(nil)[0]

	if !resolved.Type().AssignableTo(val.Elem().Type()) {
		return errors.New(fmt.Sprintf("factory result type %v is not assignable to %v", resolved.Type(), val.Elem().Type()))
	}

	val.Elem().Set(resolved)
	return nil
}

// Adds a dependency to the Container's registry using a key.
// Only one dependency can be added with a certain key.
// Dependencies are registered using factory functions that return the resolved dependency.
// Singleton behaviour is handled through the factory function's logic. If returning a pointer,
// the resolved dependency will be the same pointer.
//
// Returns an error if the dependency is already registered.
func Register[T any](c *Container, key string, factory func() T) error {
	return c.register(key, factory)
}

// Resolves a dependency and returns the resolved value.
// Returns an error if the dependency is not found.
// If the resolved value is not of type T returns an error.
func Get[T any](c *Container, key string) (T, error) {
	var result T
	err := c.get(&result, key)
	return result, err
}
