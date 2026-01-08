package main

import (
	"context"
	"fmt"
	"go-flip-life-style-products/internal/config"
	"go-flip-life-style-products/internal/contract"
	"go-flip-life-style-products/internal/deliveries/consumer"
	httpHandler "go-flip-life-style-products/internal/deliveries/http"
	workerHandler "go-flip-life-style-products/internal/deliveries/worker"
	gracefulPkg "go-flip-life-style-products/internal/pkg/graceful"
	loggerPkg "go-flip-life-style-products/internal/pkg/logger"
	uuidPkg "go-flip-life-style-products/internal/pkg/uuid"
	"os"
)

func main() {
	var (
		ctx      = context.WithValue(context.Background(), uuidPkg.CorrelationIDKey, uuidPkg.UUID())
		starters []gracefulPkg.ProcessStarter
		stoppers []gracefulPkg.ProcessStopper
	)

	cfg, err := config.New()
	if err != nil {
		loggerPkg.Fatal(ctx, fmt.Sprintf("error init config: %s", err.Error()))
		os.Exit(0)
	}

	c, stoppersContract, err := contract.New(ctx, cfg)
	if err != nil {
		loggerPkg.Fatal(ctx, fmt.Sprintf("error init contract: %s", err.Error()))
		os.Exit(0)
	}
	stoppers = append(stoppers, stoppersContract...)

	if c == nil {
		loggerPkg.Fatal(ctx, "contract is nil")
		os.Exit(0)
	}

	w := workerHandler.NewUploadWorker(c)
	w.Start(ctx, cfg.Worker.UploadWorker.Size)

	rec := consumer.NewReconciliationConsumer(c.Event.Subscribe())
	rec.Start(ctx, cfg.Consumer.ReconciliationConsumer.Size)

	httpClient := httpHandler.New(c)
	starterHTTP := httpClient.Start()
	stopperHTTP := httpClient.Stop()
	starters = append(starters, starterHTTP)
	stoppers = append(stoppers, stopperHTTP)

	gracefulPkg.StartProcessAtBackground(starters...)
	gracefulPkg.StopProcessAtBackground(ctx, cfg.App.GracefulTimeout, stoppers...)
	os.Exit(0)
}
