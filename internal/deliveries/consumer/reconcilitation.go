package consumer

import (
	"context"
	"fmt"
	loggerPkg "go-flip-life-style-products/internal/pkg/logger"
	"math/rand"
	"sync"
	"time"

	"go-flip-life-style-products/internal/models"
)

type ReconciliationConsumer struct {
	ch        <-chan models.Event
	processed sync.Map
}

func NewReconciliationConsumer(ch <-chan models.Event) *ReconciliationConsumer {
	return &ReconciliationConsumer{
		ch: ch,
	}
}

func (c *ReconciliationConsumer) Start(ctx context.Context, size int) {
	loggerPkg.Info(ctx, fmt.Sprintf("ReconciliationConsumer started with %d consumers", size))
	for i := 0; i < size; i++ {
		go c.run()
	}
}

func (c *ReconciliationConsumer) run() {
	for evt := range c.ch {
		if _, ok := c.processed.Load(evt.EventID); ok {
			continue
		}

		for i := 1; i <= 3; i++ {
			if err := c.handle(evt, i); err == nil {
				c.processed.Store(evt.EventID, true)
				break
			}
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func (c *ReconciliationConsumer) handle(evt models.Event, i int) error {
	loggerPkg.Info(evt.Ctx, fmt.Sprintf("ReconciliationConsumer handle event %s, try %d", evt.EventID, i))
	if evt.Status != models.StatusFailed {
		return nil
	}

	// simulate reconciliation
	if rng.Intn(100) < 50 {
		err := fmt.Errorf("simulated reconciliation failure %s, try %d", evt.EventID, i)
		loggerPkg.Info(evt.Ctx, err.Error())
		return err
	}

	loggerPkg.Info(evt.Ctx, fmt.Sprintf("simulated reconciliation success %s, try %d", evt.EventID, i))

	return nil
}
