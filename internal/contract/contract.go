package contract

import (
	"context"
	"go-flip-life-style-products/internal/config"
	eventPkg "go-flip-life-style-products/internal/pkg/event"
	filePkg "go-flip-life-style-products/internal/pkg/file"
	"go-flip-life-style-products/internal/pkg/graceful"
	loggerPkg "go-flip-life-style-products/internal/pkg/logger"
	queuePkg "go-flip-life-style-products/internal/pkg/queue"
)

type Contract struct {
	Cfg          *config.Configuration
	Repositories *Repositories
	Service      *Service
	File         filePkg.File
	Queue        queuePkg.Queue
	Event        eventPkg.Event
}

func New(ctx context.Context, cfg *config.Configuration) (c *Contract, stoppers []graceful.ProcessStopper, err error) {
	c = &Contract{
		Cfg: cfg,
	}

	// we can put db, redis, kafka connection here

	c.File = filePkg.NewFile()

	queue, stopperQueue := queuePkg.NewQueue(100)
	stoppers = append(stoppers, stopperQueue)
	c.Queue = queue

	event, stopperEvent := eventPkg.NewEvent(100)
	stoppers = append(stoppers, stopperEvent)
	c.Event = event

	c.Repositories, err = newRepository(c)
	if err != nil {
		loggerPkg.Fatal(ctx, "failed to init repositories")
	}

	c.Service, err = newService(c)
	if err != nil {
		loggerPkg.Fatal(ctx, "failed to init services")
	}

	return
}
