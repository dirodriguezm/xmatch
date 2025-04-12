package metadata

import (
	"log/slog"

	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/reader"
	"github.com/dirodriguezm/xmatch/service/internal/catalog_indexer/writer"
)

type IndexerActor struct {
	inbox  chan reader.ReaderResult
	outbox chan writer.WriterInput[any]
}

func New(
	inbox chan reader.ReaderResult,
	outbox chan writer.WriterInput[any],
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

func (actor *IndexerActor) receive(msg reader.ReaderResult) {
	slog.Debug("Indexer Received Message")
	if msg.Error != nil {
		actor.outbox <- writer.WriterInput[any]{
			Error: msg.Error,
			Rows:  nil,
		}
		return
	}

	objects := make([]any, len(msg.Rows))
	for i := 0; i < len(msg.Rows); i++ {
		object := msg.Rows[i]
		objects[i] = object.ToMetadata()
	}

	actor.outbox <- writer.WriterInput[any]{
		Error: nil,
		Rows:  objects,
	}
}
