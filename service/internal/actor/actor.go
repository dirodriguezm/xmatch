package actor

import (
	"context"
	"log/slog"
	"sync"
)

type Handler[T any] func(*Actor[T], Message[T])
type Stopper[T any] func(*Actor[T])

type Actor[T any] struct {
	ch       chan Message[T]
	wg       *sync.WaitGroup
	handler  Handler[T]
	stopper  Stopper[T]
	ctx      context.Context
	name     string
	stopOnce sync.Once
}

func New[T any](name string, bufferSize int, handler Handler[T], stopper Stopper[T], ctx context.Context) *Actor[T] {
	return &Actor[T]{
		name:    name,
		ch:      make(chan Message[T], bufferSize),
		wg:      &sync.WaitGroup{},
		handler: handler,
		stopper: stopper,
		ctx:     ctx,
	}
}

func (a *Actor[T]) Start() {
	a.wg.Add(1)
	go func() {
		defer func() {
			a.wg.Done()
			if err := recover(); err != nil {
				slog.Error("Actor panicked", "name", a.name)
				panic(err)
			}
		}()
		for {
			select {
			case <-a.ctx.Done():
				slog.Debug("Actor context cancellation")
				a.Stop()
			case msg, ok := <-a.ch:
				if !ok {
					slog.Debug("Actor Done")
					return
				}
				a.handler(a, msg)
			}
		}
	}()
}

// Waits for actor to finish. Stopping is idempotent so a stage stopper can
// stop its downstream actor while the pipeline also stops it explicitly.
func (a *Actor[T]) Stop() {
	a.stopOnce.Do(func() {
		close(a.ch)
		a.wg.Wait()
		if a.stopper != nil {
			a.stopper(a)
		}
	})
}

func (a *Actor[T]) Send(msg Message[T]) {
	a.ch <- msg
}
