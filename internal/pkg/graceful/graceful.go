package graceful

import (
	"context"
	"fmt"
	loggerPkg "go-flip-life-style-products/internal/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type ProcessStarter func() error

type ProcessStopper struct {
	Name string
	Stop func(ctx context.Context) error
}

func StartProcessAtBackground(ps ...ProcessStarter) {
	for _, p := range ps {
		if p != nil {
			go func(_p func() error) {
				_ = _p()
			}(p)
		}
	}
}

func StopProcessAtBackground(ctx context.Context, duration time.Duration, ps ...ProcessStopper) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigterm)

	sig := <-sigterm
	loggerPkg.Info(ctx, fmt.Sprintf("received signal %v, starting graceful shutdown", sig))

	loggerPkg.Info(ctx, fmt.Sprintf("waiting %v before stopping services", duration))
	time.Sleep(duration)

	for i, s := range ps {
		loggerPkg.Info(ctx, fmt.Sprintf("stopping %s", s.Name), log.Fields{"order": i + 1, "total": len(ps)})

		if err := s.Stop(ctx); err != nil {
			loggerPkg.Error(ctx, fmt.Sprintf("shutdown %s failed", s.Name),
				log.Fields{"error": err},
			)
		} else {
			loggerPkg.Info(ctx, fmt.Sprintf("shutdown %s successfully", s.Name),
				log.Fields{"order": i + 1, "total": len(ps)},
			)
		}
	}
}
