package writer

import "github.com/dirodriguezm/xmatch/service/internal/actor"

type Writer interface {
	Write(*actor.Actor, actor.Message)
	Stop(*actor.Actor)
}
