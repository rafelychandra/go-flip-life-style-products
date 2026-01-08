package helper

import (
	"go-flip-life-style-products/internal/config"
	mockPkgEvent "go-flip-life-style-products/internal/pkg/event/mock"
	mockPkgFile "go-flip-life-style-products/internal/pkg/file/mock"
	mockPkgQueue "go-flip-life-style-products/internal/pkg/queue/mock"
	mockRepo "go-flip-life-style-products/internal/repositories/mock"
	"go-flip-life-style-products/internal/services"
	mockService "go-flip-life-style-products/internal/services/mock"
	"testing"

	"go.uber.org/mock/gomock"
)

type TestHelper struct {
	// config
	Config *config.Configuration

	// mock ctrl
	MockCtrl *gomock.Controller

	// mock pkg
	MockFile  *mockPkgFile.MockFile
	MockEvent *mockPkgEvent.MockEvent
	MockQueue *mockPkgQueue.MockQueue

	// mock repositories
	MockTransactionRepository *mockRepo.MockTransaction

	// mock services
	MockBalanceServices    *mockService.MockBalance
	MockStatementsServices *mockService.MockStatements

	// init services
	BalanceServices     services.Balance
	StatementsServices  services.Statements
	TransactionServices services.Transaction
}

func UnitTestHelper(t *testing.T) TestHelper {
	mockCtrl := gomock.NewController(t)

	cfg := &config.Configuration{
		App:      config.App{},
		Worker:   config.Worker{},
		Consumer: config.Consumer{},
	}

	mockEventPkg := mockPkgEvent.NewMockEvent(mockCtrl)
	mockFilePkg := mockPkgFile.NewMockFile(mockCtrl)
	mockQueuePkg := mockPkgQueue.NewMockQueue(mockCtrl)

	mockTransactionRepository := mockRepo.NewMockTransaction(mockCtrl)
	mockStatementsServices := mockService.NewMockStatements(mockCtrl)
	mockBalanceServices := mockService.NewMockBalance(mockCtrl)

	balanceServices := services.NewBalance(cfg, mockTransactionRepository)
	statementsServices := services.NewStatements(cfg, mockFilePkg, mockQueuePkg)
	transactionServices := services.NewTransaction(cfg, mockTransactionRepository, mockEventPkg)

	return TestHelper{
		Config:                    cfg,
		MockCtrl:                  mockCtrl,
		MockFile:                  mockFilePkg,
		MockQueue:                 mockQueuePkg,
		MockEvent:                 mockEventPkg,
		MockTransactionRepository: mockTransactionRepository,
		MockBalanceServices:       mockBalanceServices,
		MockStatementsServices:    mockStatementsServices,
		BalanceServices:           balanceServices,
		StatementsServices:        statementsServices,
		TransactionServices:       transactionServices,
	}
}
