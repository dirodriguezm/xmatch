package assertions

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

func NotNil(object interface{}, msgAndArgs ...interface{}) {
	if object == nil {
		var msg string
		if len(msgAndArgs) > 0 {
			msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		} else {
			msg = fmt.Sprintf("expected non nil value, got %v", object)
		}
		panic(msg)
	}
}

func NotZero(object interface{}, msgAndArgs ...interface{}) {
	zero := reflect.Zero(reflect.TypeOf(object))
	if object == zero {
		var msg string
		if len(msgAndArgs) > 0 {
			msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		} else {
			msg = fmt.Sprintf("expected non zero value, got %v", object)
		}
		panic(msg)
	}
}

func HasKey(m map[string]any, key string, msgAndArgs ...interface{}) {
	_, ok := m[key]
	if !ok {
		var msg string
		if len(msgAndArgs) > 0 {
			msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		} else {
			msg = fmt.Sprintf("Expected map to have key %s", key)
		}
		panic(msg)
	}
}

func GreaterThan[T constraints.Ordered](first, second T, msgAndArgs ...interface{}) {
	if first > second {
		return
	}
	var msg string
	if len(msgAndArgs) > 0 {
		msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	} else {
		msg = fmt.Sprintf("Expected %v to be greater than %v", first, second)
	}
	panic(msg)
}
