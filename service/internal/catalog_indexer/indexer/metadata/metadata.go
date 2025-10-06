package metadata

import "github.com/dirodriguezm/xmatch/service/internal/actor"

type MetadataIndexer interface {
	Index(a *actor.Actor, msg actor.Message)
}
