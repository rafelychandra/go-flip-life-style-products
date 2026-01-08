package services

import (
	"context"
	"go-flip-life-style-products/internal/config"
	"go-flip-life-style-products/internal/models"
	"go-flip-life-style-products/internal/repositories"

	"github.com/shopspring/decimal"
)

type (
	Balance interface {
		GetBalanceByUploadID(ctx context.Context, uploadID string) (*models.Balance, error)
	}

	balance struct {
		cfg             *config.Configuration
		transactionRepo repositories.Transaction
	}
)

func NewBalance(cfg *config.Configuration, transactionRepo repositories.Transaction) Balance {
	return &balance{
		cfg:             cfg,
		transactionRepo: transactionRepo,
	}
}

func (b *balance) GetBalanceByUploadID(ctx context.Context, uploadID string) (*models.Balance, error) {
	listTransaction, _, err := b.transactionRepo.ListTransactionByUploadID(ctx, uploadID, nil)
	if err != nil {
		return nil, err
	}

	var result = models.Balance{
		UploadID: uploadID,
		Balance:  decimal.Zero,
	}

	for _, tx := range listTransaction {
		if tx.Status != models.StatusSuccess {
			continue
		}
		if tx.Type == models.TransactionTypeCredit {
			result.Balance = result.Balance.Add(tx.Amount)
		} else if tx.Type == models.TransactionTypeDebit {
			result.Balance = result.Balance.Sub(tx.Amount)
		}
	}

	return &result, nil
}
