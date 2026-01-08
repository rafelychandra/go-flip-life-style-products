package logger

import (
	"context"
	uuidPkg "go-flip-life-style-products/internal/pkg/uuid"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	log.SetLevel(log.InfoLevel)
}

func LogContext(ctx context.Context) *log.Entry {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	timeNow := time.Now().In(loc)

	fields := log.Fields{
		"service": "go-flip-life-style-products",
		"at":      timeNow.Format("2006-01-02 15:04:05"),
	}

	if ctx != nil {
		fields["correlation_id"] = uuidPkg.GetCorrelationIDFromContext(ctx)
	}

	return log.WithFields(fields)
}

type logFn func(*log.Entry, ...interface{})

func mergeFields(fields ...log.Fields) log.Fields {
	merged := log.Fields{}
	for _, f := range fields {
		for k, v := range f {
			merged[k] = v
		}
	}
	return merged
}

func logWithLevel(ctx context.Context, fn logFn, message string, fields ...log.Fields) {
	fn(LogContext(ctx).WithFields(mergeFields(fields...)), message)
}

func Info(ctx context.Context, message string, fields ...log.Fields) {
	logWithLevel(ctx, (*log.Entry).Info, message, fields...)
}

func Warn(ctx context.Context, message string, fields ...log.Fields) {
	logWithLevel(ctx, (*log.Entry).Warn, message, fields...)
}

func Error(ctx context.Context, message string, fields ...log.Fields) {
	logWithLevel(ctx, (*log.Entry).Error, message, fields...)
}

func Fatal(ctx context.Context, message string, fields ...log.Fields) {
	logWithLevel(ctx, (*log.Entry).Fatal, message, fields...)
}
