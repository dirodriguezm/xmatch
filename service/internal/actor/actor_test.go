package actor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActor(t *testing.T) {
	receivedMessages := make([]string, 0)
	handler2 := func(a *Actor, msg Message) {
		for _, row := range msg.Rows {
			receivedMessages = append(receivedMessages, row.(string))
		}
	}
	ctx := t.Context()
	actor2 := New(10, handler2, nil, nil, ctx)

	handler1 := func(a *Actor, msg Message) {
		a.Broadcast(msg)
	}
	actor1 := New(10, handler1, nil, []*Actor{actor2}, ctx)

	actor2.Start()
	actor1.Start()

	actor1.Send(Message{Rows: []any{"a", "b", "c"}, Error: nil})
	actor1.Send(Message{Rows: []any{"d", "e", "f"}, Error: nil})

	actor1.Stop()
	actor2.Stop()

	require.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, receivedMessages)
}
