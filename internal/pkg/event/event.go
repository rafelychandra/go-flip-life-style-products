package event

import (
	"context"
	"fmt"
	"go-flip-life-style-products/internal/models"
	"go-flip-life-style-products/internal/pkg/graceful"
	"sync"
)

type (
	Event interface {
		Publish(evt models.Event) error
		Subscribe() <-chan models.Event
	}

	event struct {
		eventCh chan models.Event
		stopCh  chan struct{}
		once    sync.Once
		stopped bool
		mu      sync.Mutex
	}
)

func NewEvent(size int) (Event, graceful.ProcessStopper) {
	var e = event{
		eventCh: make(chan models.Event, size),
		stopCh:  make(chan struct{}),
	}
	stopper := func(ctx context.Context) error {
		e.once.Do(func() {
			e.mu.Lock()
			e.stopped = true
			e.mu.Unlock()

			close(e.stopCh)
			close(e.eventCh)
		})
		return nil
	}

	return &e, graceful.ProcessStopper{
		Name: "event",
		Stop: stopper,
	}
}

func (e *event) Publish(evt models.Event) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.stopped {
		return fmt.Errorf("event is already stopped")
	}

	select {
	case e.eventCh <- evt:
		return nil
	case <-e.stopCh:
		return fmt.Errorf("event is already stopped")
	}
}

func (e *event) Subscribe() <-chan models.Event {
	return e.eventCh
}
