package services

import (
	"context"
	"go-flip-life-style-products/internal/config"
	"go-flip-life-style-products/internal/models"
	eventPkg "go-flip-life-style-products/internal/pkg/event"
	uuidPkg "go-flip-life-style-products/internal/pkg/uuid"
	"go-flip-life-style-products/internal/repositories"
)

type (
	Transaction interface {
		Store(ctx context.Context, tx models.Transaction) error
		GetListIssuesTransaction(ctx context.Context, req models.ReqGetListIssuesTransaction) (resp []models.Transaction, count int, limit int, err error)
	}
	transaction struct {
		cfg             *config.Configuration
		transactionRepo repositories.Transaction
		event           eventPkg.Event
	}
)

func NewTransaction(cfg *config.Configuration, transactionRepo repositories.Transaction, event eventPkg.Event) Transaction {
	return &transaction{
		cfg:             cfg,
		transactionRepo: transactionRepo,
		event:           event,
	}
}

func (s *transaction) Store(ctx context.Context, tx models.Transaction) error {
	s.transactionRepo.Add(tx)

	if tx.Status == models.StatusFailed {
		err := s.event.Publish(models.Event{
			Ctx:         ctx,
			EventID:     uuidPkg.UUID(),
			Transaction: tx,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *transaction) GetListIssuesTransaction(ctx context.Context, req models.ReqGetListIssuesTransaction) (resp []models.Transaction, count int, limit int, err error) {
	filter, err := req.PaginationFilter()
	if err != nil {
		return
	}

	limit = filter.Limit
	resp, count, err = s.transactionRepo.ListTransactionByUploadID(ctx, req.UploadID, &filter)
	if err != nil || resp == nil {
		return nil, 0, 0, nil
	}

	return

}
