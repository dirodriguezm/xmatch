package actor

import (
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActor(t *testing.T) {
	receivedMessages := make([]string, 0)
	ctx := t.Context()

	handler2 := func(_ *Actor[string], msg Message[string]) {
		receivedMessages = append(receivedMessages, msg.Rows...)
	}
	actor2 := New("sink", 10, handler2, nil, ctx)

	handler1 := func(_ *Actor[string], msg Message[string]) {
		actor2.Send(msg)
	}
	actor1 := New("source", 10, handler1, nil, ctx)

	actor2.Start()
	actor1.Start()

	actor1.Send(Message[string]{Rows: []string{"a", "b", "c"}, Error: nil})
	actor1.Send(Message[string]{Rows: []string{"d", "e", "f"}, Error: nil})

	actor1.Stop()
	actor2.Stop()

	require.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, receivedMessages)
}

func TestActor_StopDrainsBeforeStopper(t *testing.T) {
	var mu sync.Mutex
	events := make([]string, 0)
	ctx := t.Context()

	a := New(
		"draining",
		10,
		func(_ *Actor[int], _ Message[int]) {
			mu.Lock()
			defer mu.Unlock()
			events = append(events, "handler")
		},
		func(_ *Actor[int]) {
			mu.Lock()
			defer mu.Unlock()
			events = append(events, "stopper")
		},
		ctx,
	)
	a.Start()

	for i := range 3 {
		a.Send(Message[int]{Rows: []int{i}})
	}
	a.Stop()

	require.Equal(t, []string{"handler", "handler", "handler", "stopper"}, events)
}

func TestActor_StopChainsToDownstreamOnce(t *testing.T) {
	var mu sync.Mutex
	events := make([]string, 0)
	ctx := t.Context()

	sink := New(
		"writer",
		10,
		func(_ *Actor[int], _ Message[int]) {},
		func(_ *Actor[int]) {
			mu.Lock()
			defer mu.Unlock()
			events = append(events, "writer")
		},
		ctx,
	)
	indexer := New(
		"indexer",
		10,
		func(_ *Actor[int], msg Message[int]) { sink.Send(msg) },
		func(_ *Actor[int]) {
			mu.Lock()
			events = append(events, "indexer")
			mu.Unlock()
			sink.Stop()
		},
		ctx,
	)

	sink.Start()
	indexer.Start()
	indexer.Send(Message[int]{Rows: []int{1}})

	indexer.Stop()
	// A chained stop followed by an explicit stop must not stop twice.
	indexer.Stop()
	sink.Stop()

	require.Equal(t, []string{"indexer", "writer"}, events)
}

func TestActor_PanicPropagates(t *testing.T) {
	if os.Getenv("ACTOR_PANIC_HELPER") == "1" {
		a := New(
			"panicking",
			1,
			func(_ *Actor[int], _ Message[int]) { panic("actor boom") },
			nil,
			t.Context(),
		)
		a.Start()
		a.Send(Message[int]{Rows: []int{1}})
		a.Stop()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestActor_PanicPropagates$")
	cmd.Env = append(os.Environ(), "ACTOR_PANIC_HELPER=1")
	output, err := cmd.CombinedOutput()

	require.Error(t, err, "the panicking actor should crash its process")
	assert.Contains(t, string(output), "actor boom")
}
