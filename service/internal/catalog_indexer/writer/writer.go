package writer

import "github.com/dirodriguezm/xmatch/service/internal/actor"

type Writer[T any] interface {
	Write(*actor.Actor[T], actor.Message[T])
	Stop(*actor.Actor[T])
}
