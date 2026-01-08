package models

import "context"

type Event struct {
	Ctx     context.Context
	EventID string
	Transaction
}
