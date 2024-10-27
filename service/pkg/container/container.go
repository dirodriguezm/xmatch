package container

import (
	"fmt"
	"reflect"
)

type binding struct {
	factory  any
	instance any
}

func (b *binding) make(ctr Container, names []string) (any, error) {
	if b.instance != nil {
		return b.instance, nil
	}
	instance, err := ctr.invoke(b.factory, names)
	if err != nil {
		return nil, err
	}
	b.instance = instance
	return instance, nil
}

type Container map[reflect.Type]map[string]*binding

func NewContainer() Container {
	return make(Container)
}

func (ctr Container) Reset() {
	for d := range ctr {
		delete(ctr, d)
	}
}

func (ctr Container) Register(name string, factory any) error {
	if err := ctr.validateFactoryFunction(factory); err != nil {
		return err
	}
	returnType := reflect.TypeOf(factory).Out(0)
	if _, exist := ctr[returnType]; !exist {
		ctr[returnType] = make(map[string]*binding)
	}
	ctr[returnType][name] = &binding{factory: factory, instance: nil} // instance is set on resolve
	return nil
}

func (ctr Container) Resolve(name string, emptyValue any) error {
	valueType := reflect.TypeOf(emptyValue)
	if valueType == nil {
		return fmt.Errorf("value to resolve to can't be nil")
	}
	if valueType.Kind() != reflect.Ptr {
		return fmt.Errorf("value to resolve to must be a pointer")
	}
	bind, exist := ctr[valueType.Elem()][name]
	if !exist {
		return fmt.Errorf("dependency %s not found in container", name)
	}
	instance, err := bind.make(ctr, []string{})
	if err != nil {
		return fmt.Errorf("Could not resolve dependency %s. Error: %w", name, err)
	}
	reflect.ValueOf(emptyValue).Elem().Set(reflect.ValueOf(instance))
	return nil
}

func (ctr Container) ResolveWithBinds(name string, emptyValue any, names []string) error {
	valueType := reflect.TypeOf(emptyValue)
	if valueType == nil {
		return fmt.Errorf("value to resolve to can't be nil")
	}
	if valueType.Kind() != reflect.Ptr {
		return fmt.Errorf("value to resolve to must be a pointer")
	}
	bind, exist := ctr[valueType.Elem()][name]
	if !exist {
		return fmt.Errorf("dependency %s not found in container", name)
	}
	instance, err := bind.make(ctr, names)
	if err != nil {
		return fmt.Errorf("Could not resolve dependency %s. Error: %w", name, err)
	}
	reflect.ValueOf(emptyValue).Elem().Set(reflect.ValueOf(instance))
	return nil
}

func (ctr *Container) validateFactoryFunction(factory any) error {
	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		return fmt.Errorf("factory argument must be a function")
	}
	numRets := factoryType.NumOut()
	if numRets == 0 || numRets > 2 {
		return fmt.Errorf("factory function can only return an interface or an interface and error")
	}

	returnType := factoryType.Out(0)
	for i := 0; i < factoryType.NumIn(); i++ {
		if factoryType.In(i) == returnType {
			return fmt.Errorf("can't depend on the same type it returns")
		}
	}
	return nil
}

func (ctr Container) invoke(factory any, names []string) (any, error) {
	arguments, err := ctr.arguments(factory, names)
	if err != nil {
		return nil, err
	}

	values := reflect.ValueOf(factory).Call(arguments)
	fmt.Printf("VALUES %v", values)
	if len(values) == 2 && values[1].CanInterface() {
		if err, ok := values[1].Interface().(error); ok {
			return values[0].Interface(), err
		}
	}
	return values[0].Interface(), nil
}

type NamedDependency struct {
	Name    string
	Factory any
}

func (c Container) arguments(factory any, names []string) ([]reflect.Value, error) {
	factoryType := reflect.TypeOf(factory)
	argumentsCount := factoryType.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := factoryType.In(i)
		fmt.Printf("Abstraction %v \t Name %s\n", abstraction, names[i])
		fmt.Printf("Container: %+v\n", c)
		if bind, exist := c[abstraction][names[i]]; exist {
			instance, err := bind.make(c, names)
			if err != nil {
				return nil, err
			}
			arguments[i] = reflect.ValueOf(instance)
		} else {
			return nil, fmt.Errorf("no instance found for: %s", abstraction)
		}
	}

	return arguments, nil
}
