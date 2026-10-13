package actor

type Message[T any] struct {
	Rows  []T
	Error error
}
