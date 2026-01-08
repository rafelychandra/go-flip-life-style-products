package repositories

import (
	"context"
	"fmt"
	"go-flip-life-style-products/internal/models"
	loggerPkg "go-flip-life-style-products/internal/pkg/logger"
	"sort"
	"sync"

	log "github.com/sirupsen/logrus"
)

type (
	Transaction interface {
		Add(tx models.Transaction)
		ListTransactionByUploadID(ctx context.Context, uploadID string, filter *models.FilterTransaction) ([]models.Transaction, int, error)
	}

	transaction struct {
		mu   sync.Mutex
		data map[string][]models.Transaction
	}
)

func NewTransaction() Transaction {
	return &transaction{
		data: make(map[string][]models.Transaction),
	}
}

func (s *transaction) Add(tx models.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[tx.UploadID] = append(s.data[tx.UploadID], tx)
}

func (s *transaction) ListTransactionByUploadID(ctx context.Context, uploadID string, filter *models.FilterTransaction) ([]models.Transaction, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, ok := s.data[uploadID]
	if !ok {
		loggerPkg.Error(ctx, "data not found", log.Fields{"upload_id": uploadID})
		return nil, 0, fmt.Errorf("data not found, uploadID: %s", uploadID)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Timestamp > data[j].Timestamp
	})

	if filter != nil {
		return s.listTransactionWithFilter(data, filter)
	}

	return data, len(data), nil
}

func (s *transaction) listTransactionWithFilter(data []models.Transaction, filter *models.FilterTransaction) ([]models.Transaction, int, error) {
	filtered := make([]models.Transaction, 0, len(data))
	for _, tx := range data {
		if len(filter.Status) > 0 && !containsStatus(filter.Status, tx.Status) {
			continue
		}

		if len(filter.TransactionType) > 0 && !containsType(filter.TransactionType, tx.Type) {
			continue
		}
		filtered = append(filtered, tx)
	}

	paged := make([]models.Transaction, 0, filter.Limit)
	for _, tx := range filtered {
		if filter.NextCursor > 0 && tx.Timestamp >= filter.NextCursor {
			continue
		}
		if filter.PrevCursor > 0 && tx.Timestamp <= filter.PrevCursor {
			continue
		}

		if len(paged) < filter.Limit {
			paged = append(paged, tx)
		}
	}

	if filter.PrevCursor > 0 {
		for i, j := 0, len(paged)-1; i < j; i, j = i+1, j-1 {
			paged[i], paged[j] = paged[j], paged[i]
		}
	}

	return paged, len(filtered), nil
}

func containsStatus(list []models.Status, v models.Status) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func containsType(list []models.TransactionType, v models.TransactionType) bool {
	for _, t := range list {
		if t == v {
			return true
		}
	}
	return false
}
