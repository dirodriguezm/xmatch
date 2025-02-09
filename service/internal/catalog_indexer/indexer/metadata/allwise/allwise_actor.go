package allwise_metadata

import (
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/indexer"
	"github.com/dirodriguezm/xmatch/service/internal/repository"
)

type IndexerActor struct {
	inbox  chan indexer.ReaderResult
	outbox chan indexer.WriterInput[repository.AllwiseMetadata]
}

func New(
	inbox chan indexer.ReaderResult,
	outbox chan indexer.WriterInput[repository.AllwiseMetadata],
) *IndexerActor {
	slog.Debug("Creating new Indexer")
	return &IndexerActor{
		inbox:  inbox,
		outbox: outbox,
	}

}

// Starts the Actor goroutine
//
// Listen for messages in the inbox channel and process them
// sending the results to the outbox channel
func (actor *IndexerActor) Start() {
	slog.Debug("Starting Indexer")
	go func() {
		defer func() {
			close(actor.outbox)
			slog.Debug("Closing Indexer")
		}()
		for msg := range actor.inbox {
			actor.receive(msg)
		}
	}()
}

func (actor *IndexerActor) receive(msg indexer.ReaderResult) {
	slog.Debug("Indexer Received Message")
	if msg.Error != nil {
		actor.outbox <- indexer.WriterInput[repository.AllwiseMetadata]{
			Error: msg.Error,
			Rows:  nil,
		}
		return
	}
	objects := make([]repository.AllwiseMetadata, len(msg.Rows))
	for i := 0; i < len(msg.Rows); i++ {
		object := msg.Rows[i].(*repository.AllwiseInputSchema)
		objects[i] = object.ToMetadata()
	}
	actor.outbox <- indexer.WriterInput[repository.AllwiseMetadata]{
		Error: nil,
		Rows:  objects,
	}
}
