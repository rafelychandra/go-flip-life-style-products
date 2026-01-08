package services_test

import (
	"context"
	"fmt"
	"go-flip-life-style-products/internal/models"
	helperPkg "go-flip-life-style-products/internal/pkg/tester"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStore(t *testing.T) {
	var (
		unitTestHelper = helperPkg.UnitTestHelper(t)
		ctx            = context.Background()
	)

	type args struct {
		ctx context.Context
		tx  models.Transaction
	}

	tests := []struct {
		name          string
		args          args
		wantErr       bool
		expectedError error
		doMockService func(helperPkg.TestHelper, args)
	}{
		{
			name: "success store SUCCESS data",
			args: args{
				ctx: ctx,
				tx: models.Transaction{
					UploadID:     "123456",
					Timestamp:    123456789,
					Counterparty: "JOHN DOE",
					Type:         "DEBIT",
					Amount:       decimal.RequireFromString("10"),
					Status:       "SUCCESS",
					Description:  "description",
				},
			},
			wantErr:       false,
			expectedError: nil,
			doMockService: func(helper helperPkg.TestHelper, args args) {
				helper.MockTransactionRepository.EXPECT().Add(args.tx)
			},
		},
		{
			name: "success store FAILED data",
			args: args{
				ctx: ctx,
				tx: models.Transaction{
					UploadID:     "123456",
					Timestamp:    123456789,
					Counterparty: "JOHN DOE",
					Type:         "DEBIT",
					Amount:       decimal.RequireFromString("10"),
					Status:       "FAILED",
					Description:  "description",
				},
			},
			wantErr:       false,
			expectedError: nil,
			doMockService: func(helper helperPkg.TestHelper, args args) {
				helper.MockTransactionRepository.EXPECT().Add(args.tx)
				helper.MockEvent.EXPECT().Publish(gomock.Any()).Return(nil)
			},
		},
		{
			name: "success store PENDING data",
			args: args{
				ctx: ctx,
				tx: models.Transaction{
					UploadID:     "123456",
					Timestamp:    123456789,
					Counterparty: "JOHN DOE",
					Type:         "DEBIT",
					Amount:       decimal.RequireFromString("10"),
					Status:       "PENDING",
					Description:  "description",
				},
			},
			wantErr:       false,
			expectedError: nil,
			doMockService: func(helper helperPkg.TestHelper, args args) {
				helper.MockTransactionRepository.EXPECT().Add(args.tx)
			},
		},
		{
			name: "error publish FAILED data",
			args: args{
				ctx: ctx,
				tx: models.Transaction{
					UploadID:     "123456",
					Timestamp:    123456789,
					Counterparty: "JOHN DOE",
					Type:         "DEBIT",
					Amount:       decimal.RequireFromString("10"),
					Status:       "FAILED",
					Description:  "description",
				},
			},
			wantErr:       true,
			expectedError: fmt.Errorf("got an error"),
			doMockService: func(helper helperPkg.TestHelper, args args) {
				helper.MockTransactionRepository.EXPECT().Add(args.tx)
				helper.MockEvent.EXPECT().Publish(gomock.Any()).Return(fmt.Errorf("got an error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.doMockService(unitTestHelper, tt.args)

			err := unitTestHelper.TransactionServices.Store(tt.args.ctx, tt.args.tx)
			if (err != nil) == tt.wantErr {
				assert.Equal(t, tt.expectedError, err)
			}

			if !tt.wantErr {
				assert.NoError(t, tt.expectedError, err)
			}
		})
	}
}

func TestGetListIssuesTransaction(t *testing.T) {
	var (
		unitTestHelper = helperPkg.UnitTestHelper(t)
		ctx            = context.Background()
	)

	type args struct {
		ctx context.Context
		req models.ReqGetListIssuesTransaction
	}

	type expectedOutput struct {
		resp  []models.Transaction
		count int
		limit int
	}

	tests := []struct {
		name           string
		args           args
		wantErr        bool
		expectedError  error
		expectedOutput expectedOutput
		doMockService  func(helperPkg.TestHelper, args)
	}{
		{
			name: "success but no data",
			args: args{
				ctx: ctx,
				req: models.ReqGetListIssuesTransaction{
					UploadID: "123456",
				},
			},
			wantErr:       false,
			expectedError: nil,
			expectedOutput: expectedOutput{
				resp:  nil,
				count: 0,
				limit: 0,
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.req.UploadID, gomock.Any()).Return(nil, 0, nil)
			},
		},
		{
			name: "success with no filter (default status FAILED, PENDING)",
			args: args{
				ctx: ctx,
				req: models.ReqGetListIssuesTransaction{
					UploadID: "123456",
				},
			},
			wantErr:       false,
			expectedError: nil,
			expectedOutput: expectedOutput{
				resp: []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10"),
						Status:       "FAILED",
						Description:  "description",
					},
				},
				count: 1,
				limit: 11,
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {
				responseList := []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10"),
						Status:       "FAILED",
						Description:  "description",
					},
				}
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.req.UploadID, gomock.Any()).Return(responseList, 1, nil)
			},
		},
		{
			name: "success with filter status = SUCCESS",
			args: args{
				ctx: ctx,
				req: models.ReqGetListIssuesTransaction{
					UploadID: "123456",
					Status:   "SUCCESS",
				},
			},
			wantErr:       false,
			expectedError: nil,
			expectedOutput: expectedOutput{
				resp: []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10000"),
						Status:       "SUCCESS",
						Description:  "description",
					},
				},
				count: 1,
				limit: 11,
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {
				responseList := []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10000"),
						Status:       "SUCCESS",
						Description:  "description",
					},
				}
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.req.UploadID, gomock.Any()).Return(responseList, 1, nil)
			},
		},
		{
			name: "success with limit 1 and default status (FAILED,PENDING)",
			args: args{
				ctx: ctx,
				req: models.ReqGetListIssuesTransaction{
					UploadID: "123456",
					Limit:    1,
				},
			},
			wantErr:       false,
			expectedError: nil,
			expectedOutput: expectedOutput{
				resp: []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10000"),
						Status:       "FAILED",
						Description:  "description",
					},
				},
				count: 1,
				limit: 2,
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {
				responseList := []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10000"),
						Status:       "FAILED",
						Description:  "description",
					},
				}
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.req.UploadID, gomock.Any()).Return(responseList, 1, nil)
			},
		},
		{
			name: "success with next cursor",
			args: args{
				ctx: ctx,
				req: models.ReqGetListIssuesTransaction{
					UploadID:   "123456",
					NextCursor: 123456789,
				},
			},
			wantErr:       false,
			expectedError: nil,
			expectedOutput: expectedOutput{
				resp: []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10000"),
						Status:       "FAILED",
						Description:  "description",
					},
				},
				count: 1,
				limit: 11,
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {
				responseList := []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10000"),
						Status:       "FAILED",
						Description:  "description",
					},
				}
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.req.UploadID, gomock.Any()).Return(responseList, 1, nil)
			},
		},
		{
			name: "success with prev cursor",
			args: args{
				ctx: ctx,
				req: models.ReqGetListIssuesTransaction{
					UploadID:   "123456",
					PrevCursor: 123456789,
				},
			},
			wantErr:       false,
			expectedError: nil,
			expectedOutput: expectedOutput{
				resp: []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10000"),
						Status:       "FAILED",
						Description:  "description",
					},
				},
				count: 1,
				limit: 11,
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {
				responseList := []models.Transaction{
					{
						UploadID:     "123456",
						Timestamp:    123456789,
						Counterparty: "JOHN DOE",
						Type:         "DEBIT",
						Amount:       decimal.RequireFromString("10000"),
						Status:       "FAILED",
						Description:  "description",
					},
				}
				helper.MockTransactionRepository.EXPECT().ListTransactionByUploadID(args.ctx, args.req.UploadID, gomock.Any()).Return(responseList, 1, nil)
			},
		},
		{
			name: "error when next cursor and prev cursor filled",
			args: args{
				ctx: ctx,
				req: models.ReqGetListIssuesTransaction{
					UploadID:   "123456",
					NextCursor: 123456789,
					PrevCursor: 123456789,
				},
			},
			wantErr:       true,
			expectedError: fmt.Errorf("next_cursor and prev_cursor cannot be used together"),
			expectedOutput: expectedOutput{
				resp:  nil,
				count: 0,
				limit: 0,
			},
			doMockService: func(helper helperPkg.TestHelper, args args) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.doMockService(unitTestHelper, tt.args)

			res, count, limit, err := unitTestHelper.TransactionServices.GetListIssuesTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) == tt.wantErr {
				assert.Equal(t, tt.expectedError, err)
			}

			if !tt.wantErr {
				assert.NoError(t, tt.expectedError, err)
			}

			assert.Equal(t, tt.expectedOutput.resp, res)
			assert.Equal(t, tt.expectedOutput.count, count)
			assert.Equal(t, tt.expectedOutput.limit, limit)
		})
	}
}
