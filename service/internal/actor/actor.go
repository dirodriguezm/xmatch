package actor

import (
	"context"
	"log/slog"
	"sync"
)

type Handler func(*Actor, Message)
type Stopper func(*Actor)

type Actor struct {
	ch        chan Message
	wg        *sync.WaitGroup
	handler   Handler
	stopper   Stopper
	receivers []*Actor
	ctx       context.Context
	name      string
}

func New(name string, bufferSize int, handler Handler, stopper Stopper, receivers []*Actor, ctx context.Context) *Actor {
	return &Actor{
		name:      name,
		ch:        make(chan Message, bufferSize),
		wg:        &sync.WaitGroup{},
		handler:   handler,
		stopper:   stopper,
		receivers: receivers,
		ctx:       ctx,
	}
}

func (a *Actor) Start() {
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

// Waits for actor to finish
func (a *Actor) Stop() {
	close(a.ch)
	a.wg.Wait()
	if a.stopper != nil {
		a.stopper(a)
	}
}

func (a *Actor) Send(msg Message) {
	a.ch <- msg
}

func (a *Actor) Broadcast(msg Message) {
	for _, receiver := range a.receivers {
		receiver.Send(msg)
	}
}
