package utils

import (
	"reflect"
)

// Merges two structs. b takes precedence, but a provides defaults.
//
// Warning: This function was vibe-coded
func MergeStructs[T any](a, b T) (T, error) {
	aValue := reflect.ValueOf(a)
	bValue := reflect.ValueOf(b)

	// Create a new instance
	result := reflect.New(aValue.Type()).Elem()

	for i := 0; i < aValue.NumField(); i++ {
		aField := aValue.Field(i)
		bField := bValue.Field(i)
		resultField := result.Field(i)

		if !resultField.CanSet() {
			continue
		}

		if !bField.IsZero() {
			if bField.Kind() == reflect.Struct {
				mergedNested, err := MergeStructs(aField.Interface(), bField.Interface())
				if err != nil {
					var zero T
					return zero, err
				}
				resultField.Set(reflect.ValueOf(mergedNested))
			} else {
				resultField.Set(bField)
			}
		} else {
			resultField.Set(aField)
		}
	}

	return result.Interface().(T), nil
}
