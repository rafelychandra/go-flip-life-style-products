package services_test

import (
	"context"
	"fmt"
	"go-flip-life-style-products/internal/models"
	helperPkg "go-flip-life-style-products/internal/pkg/tester"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetBalanceByUploadID(t *testing.T) {
	var (
		unitTestHelper = helperPkg.UnitTestHelper(t)
		ctx            = context.Background()
	)

	type args struct {
		ctx      context.Context
		uploadID string
	}

	tests := []struct {
		name           string
		args           args
		wantErr        bool
		expectedError  error
		expectedOutput *models.Balance
		doMockService  func(helperPkg.TestHelper, args)
	}{
		{
			name: "success but the data is only CREDIT",
			args: args{
				ctx:      ctx,
				uploadID: "123456789",
			},
			wantErr:       false,
			expectedError: nil,
			expectedOutput: &models.Balance{
				UploadID: "123456789",
				Balance:  decimal.RequireFromString("10000"),
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {
				var outputListTransaction = []models.Transaction{
					{
						UploadID:     "123456789",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "CREDIT",
						Amount:       decimal.Decimal{}.Add(decimal.NewFromFloat(10000)),
						Status:       "SUCCESS",
						Description:  "description",
					},
					{
						UploadID:     "123456789",
						Timestamp:    987654321,
						Counterparty: "DOE JOHN",
						Type:         "DEBIT",
						Amount:       decimal.Decimal{}.Add(decimal.NewFromFloat(10000)),
						Status:       "FAILED",
						Description:  "description",
					},
				}
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.uploadID, nil).Return(outputListTransaction, 0, nil)
			},
		},
		{
			name: "success calculated the balance data",
			args: args{
				ctx:      ctx,
				uploadID: "123456789",
			},
			wantErr:       false,
			expectedError: nil,
			expectedOutput: &models.Balance{
				UploadID: "123456789",
				Balance:  decimal.RequireFromString("5000"),
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {
				var outputListTransaction = []models.Transaction{
					{
						UploadID:     "123456789",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "CREDIT",
						Amount:       decimal.Decimal{}.Add(decimal.NewFromFloat(10000)),
						Status:       "SUCCESS",
						Description:  "description",
					},
					{
						UploadID:     "123456789",
						Timestamp:    987654321,
						Counterparty: "DOE JOHN",
						Type:         "DEBIT",
						Amount:       decimal.Decimal{}.Add(decimal.NewFromFloat(5000)),
						Status:       "SUCCESS",
						Description:  "description",
					},
				}
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.uploadID, nil).Return(outputListTransaction, 0, nil)
			},
		},
		{
			name: "success calculated the balance data but the balance is negative",
			args: args{
				ctx:      ctx,
				uploadID: "123456789",
			},
			wantErr:       false,
			expectedError: nil,
			expectedOutput: &models.Balance{
				UploadID: "123456789",
				Balance:  decimal.RequireFromString("-5000"),
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {
				var outputListTransaction = []models.Transaction{
					{
						UploadID:     "123456789",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "CREDIT",
						Amount:       decimal.Decimal{}.Add(decimal.NewFromFloat(5000)),
						Status:       "SUCCESS",
						Description:  "description",
					},
					{
						UploadID:     "123456789",
						Timestamp:    987654321,
						Counterparty: "DOE JOHN",
						Type:         "DEBIT",
						Amount:       decimal.Decimal{}.Add(decimal.NewFromFloat(10000)),
						Status:       "SUCCESS",
						Description:  "description",
					},
				}
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.uploadID, nil).Return(outputListTransaction, 0, nil)
			},
		},
		{
			name: "error get list transactions for calculate balance",
			args: args{
				ctx:      ctx,
				uploadID: "123456789",
			},
			wantErr:        true,
			expectedError:  fmt.Errorf("got an error"),
			expectedOutput: nil,
			doMockService: func(helper helperPkg.TestHelper, args args) {
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.uploadID, nil).Return(nil, 0, fmt.Errorf("got an error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.doMockService(unitTestHelper, tt.args)

			res, err := unitTestHelper.BalanceServices.GetBalanceByUploadID(tt.args.ctx, tt.args.uploadID)
			if (err != nil) == tt.wantErr {
				assert.Equal(t, tt.expectedError, err)
			}

			if !tt.wantErr {
				assert.NoError(t, tt.expectedError, err)
			}

			assert.Equal(t, tt.expectedOutput, res)
		})
	}
}
