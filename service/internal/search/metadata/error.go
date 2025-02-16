package metadata

import "fmt"

type ValidationError struct {
	Field  string
	Reason string
	Value  string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf(
		"Could not parse field %s with value %s: %s",
		e.Field,
		e.Value,
		e.Reason,
	)
}

type ArgumentError struct {
	Reason string
	Value  string
	Name   string
}

func (e ArgumentError) Error() string {
	return fmt.Sprintf(
		"Could not parse argument %s with value %s: %s",
		e.Name,
		e.Value,
		e.Reason,
	)
}
