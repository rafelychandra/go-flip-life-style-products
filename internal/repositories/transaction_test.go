package repositories

import (
	"context"
	"fmt"
	"go-flip-life-style-products/internal/models"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	store := &transaction{
		data: make(map[string][]models.Transaction),
	}

	const workers = 100
	const uploadID = "upload-1"

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			store.Add(models.Transaction{
				UploadID:     uploadID,
				Timestamp:    123456789,
				Counterparty: "JOHN DOE",
				Type:         "CREDIT",
				Amount:       decimal.RequireFromString("10000"),
				Status:       "SUCCESS",
				Description:  "description",
			})
		}(i)
	}

	wg.Wait()

	assert.Len(t, store.data[uploadID], workers)
}

func TestTransactionListTransactionWithoutFilter(t *testing.T) {
	ctx := context.Background()

	store := &transaction{
		data: map[string][]models.Transaction{
			"upload-1": {
				{UploadID: "upload-1", Timestamp: 100},
				{UploadID: "upload-1", Timestamp: 300},
				{UploadID: "upload-1", Timestamp: 200},
			},
		},
	}

	result, total, err := store.ListTransactionByUploadID(ctx, "upload-1", nil)

	assert.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Len(t, result, 3)

	assert.Equal(t, int64(300), result[0].Timestamp)
	assert.Equal(t, int64(200), result[1].Timestamp)
	assert.Equal(t, int64(100), result[2].Timestamp)
}

func TestTransactionListTransactionWithoutFilterError(t *testing.T) {
	ctx := context.Background()

	store := &transaction{
		data: map[string][]models.Transaction{
			"upload-2": {
				{UploadID: "upload-1", Timestamp: 100},
				{UploadID: "upload-1", Timestamp: 300},
				{UploadID: "upload-1", Timestamp: 200},
			},
		},
	}

	result, total, err := store.ListTransactionByUploadID(ctx, "upload-1", nil)

	assert.Equal(t, fmt.Errorf("data not found, uploadID: %s", "upload-1"), err)
	assert.Equal(t, 0, total)
	assert.Len(t, result, 0)
}

func TestTransactionListTransactionWithFilterTransactionTypeAndPagination(t *testing.T) {
	ctx := context.Background()

	store := &transaction{
		data: map[string][]models.Transaction{
			"upload-1": {
				{UploadID: "upload-1", Timestamp: 100, Type: models.TransactionTypeCredit},
				{UploadID: "upload-1", Timestamp: 200, Type: models.TransactionTypeDebit},
			},
		},
	}

	filter := &models.FilterTransaction{
		Limit:           10,
		NextCursor:      0,
		PrevCursor:      0,
		TransactionType: []models.TransactionType{models.TransactionTypeCredit},
	}

	result, total, err := store.ListTransactionByUploadID(ctx, "upload-1", filter)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, result, 1)
}

func TestTransactionListTransactionWithFilterStatusAndPagination(t *testing.T) {
	ctx := context.Background()

	store := &transaction{
		data: map[string][]models.Transaction{
			"upload-1": {
				{UploadID: "upload-1", Timestamp: 100, Type: models.TransactionTypeCredit, Status: models.StatusSuccess},
				{UploadID: "upload-1", Timestamp: 200, Type: models.TransactionTypeDebit, Status: models.StatusFailed},
			},
		},
	}

	filter := &models.FilterTransaction{
		Limit:      10,
		NextCursor: 0,
		PrevCursor: 0,
		Status:     []models.Status{models.StatusSuccess},
	}

	result, total, err := store.ListTransactionByUploadID(ctx, "upload-1", filter)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, result, 1)
}

func TestTransactionListTransactionWithFilterNextCursorAndPagination(t *testing.T) {
	ctx := context.Background()

	store := &transaction{
		data: map[string][]models.Transaction{
			"upload-1": {
				{UploadID: "upload-1", Timestamp: 100, Type: models.TransactionTypeCredit, Status: models.StatusPending},
				{UploadID: "upload-1", Timestamp: 200, Type: models.TransactionTypeDebit, Status: models.StatusFailed},
			},
		},
	}

	filter := &models.FilterTransaction{
		Limit:           1,
		NextCursor:      200,
		PrevCursor:      0,
		Status:          []models.Status{models.StatusPending, models.StatusFailed},
		TransactionType: []models.TransactionType{models.TransactionTypeDebit, models.TransactionTypeCredit},
	}

	result, total, err := store.ListTransactionByUploadID(ctx, "upload-1", filter)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 1)
}

func TestTransactionListTransactionWithFilterPrevCursorAndPagination(t *testing.T) {
	ctx := context.Background()

	store := &transaction{
		data: map[string][]models.Transaction{
			"upload-1": {
				{UploadID: "upload-1", Timestamp: 100, Type: models.TransactionTypeCredit, Status: models.StatusPending},
				{UploadID: "upload-1", Timestamp: 200, Type: models.TransactionTypeDebit, Status: models.StatusFailed},
			},
		},
	}

	filter := &models.FilterTransaction{
		Limit:           1,
		NextCursor:      0,
		PrevCursor:      100,
		Status:          []models.Status{models.StatusPending, models.StatusFailed},
		TransactionType: []models.TransactionType{models.TransactionTypeDebit, models.TransactionTypeCredit},
	}

	result, total, err := store.ListTransactionByUploadID(ctx, "upload-1", filter)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 1)
}
