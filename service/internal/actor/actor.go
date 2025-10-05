package actor

import "sync"

type Handler func(*Actor, Message)

type Actor struct {
	ch        chan Message
	wg        *sync.WaitGroup
	handler   Handler
	receivers []*Actor
}

func New(bufferSize int, handler Handler, receivers []*Actor) *Actor {
	return &Actor{
		ch:        make(chan Message, bufferSize),
		wg:        &sync.WaitGroup{},
		handler:   handler,
		receivers: receivers,
	}
}

func (a *Actor) Start() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for msg := range a.ch {
			a.handler(a, msg)
		}
	}()
}

func (a *Actor) Stop() {
	close(a.ch)
	a.wg.Wait()
}

func (a *Actor) Send(msg Message) {
	a.ch <- msg
}
