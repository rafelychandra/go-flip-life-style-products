package uuid

import (
	"context"

	"github.com/google/uuid"
)

const (
	CorrelationIDKey = "CorrelationID"
)

func UUID() string {
	return uuid.New().String()
}

func GetCorrelationIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(CorrelationIDKey).(string)
	if !ok {
		return ""
	}

	return id
}
