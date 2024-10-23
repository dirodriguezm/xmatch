package assertions

import (
	"fmt"
	"reflect"
)

func NotNil(object interface{}, msgAndArgs ...interface{}) {
	if object == nil {
		msg := fmt.Sprintf("expected non nil value, got %v", object)
		panic(msg)
	}
}

func NotZero(object interface{}, msgAndArgs ...interface{}) {
	zero := reflect.Zero(reflect.TypeOf(object))
	if object == zero {
		msg := fmt.Sprintf("expected non zero value, got %v", object)
		panic(msg)
	}
}
