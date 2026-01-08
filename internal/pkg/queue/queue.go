package queue

import (
	"context"
	"fmt"
	"go-flip-life-style-products/internal/models"
	"go-flip-life-style-products/internal/pkg/graceful"
	"sync"
)

type (
	Queue interface {
		Enqueue(job models.UploadQueue) error
		Channel() <-chan models.UploadQueue
	}

	queue struct {
		jobCh   chan models.UploadQueue
		stopCh  chan struct{}
		once    sync.Once
		stopped bool
		mu      sync.Mutex
	}
)

func NewQueue(buffer int) (Queue, graceful.ProcessStopper) {
	var q = queue{
		jobCh:  make(chan models.UploadQueue, buffer),
		stopCh: make(chan struct{}),
	}

	stopper := func(ctx context.Context) error {
		q.once.Do(func() {
			q.mu.Lock()
			q.stopped = true
			q.mu.Unlock()

			close(q.stopCh)
			close(q.jobCh)
		})
		return nil
	}

	return &q, graceful.ProcessStopper{
		Name: "queue",
		Stop: stopper,
	}
}
func (q *queue) Enqueue(job models.UploadQueue) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.stopped {
		return fmt.Errorf("queue is already stopped")
	}

	select {
	case q.jobCh <- job:
		return nil
	case <-q.stopCh:
		return fmt.Errorf("queue is already stopped")
	}
}

func (q *queue) Channel() <-chan models.UploadQueue {
	return q.jobCh
}
