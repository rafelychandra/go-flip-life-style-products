package models

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusSuccess Status = "SUCCESS"
	StatusFailed  Status = "FAILED"
	StatusPending Status = "PENDING"
)

type TransactionType string

const (
	TransactionTypeCredit TransactionType = "CREDIT"
	TransactionTypeDebit  TransactionType = "DEBIT"
)

type (
	Transaction struct {
		UploadID     string          `json:"uploadID"`
		Timestamp    int64           `json:"timestamp"`
		Counterparty string          `json:"counterParty"`
		Type         TransactionType `json:"type"`
		Amount       decimal.Decimal `json:"amount"`
		Status       Status          `json:"status"`
		Description  string          `json:"description"`
	}
	ReqGetListIssuesTransaction struct {
		UploadID        string `query:"upload_id"`
		Limit           int    `query:"limit"`
		NextCursor      int64  `query:"next_cursor"`
		PrevCursor      int64  `query:"prev_cursor"`
		Status          string `query:"status"`
		TransactionType string `query:"transaction_type"`
	}
	FilterTransaction struct {
		Limit           int
		NextCursor      int64
		PrevCursor      int64
		Status          []Status
		TransactionType []TransactionType
	}
)

func (req ReqGetListIssuesTransaction) PaginationFilter() (FilterTransaction, error) {
	if req.Limit == 0 {
		req.Limit = 10
	}

	// using over-fetch limit to check next page exists or not
	req.Limit += 1

	if req.NextCursor > 0 && req.PrevCursor > 0 {
		return FilterTransaction{}, fmt.Errorf("next_cursor and prev_cursor cannot be used together")
	}

	var listStatus []Status
	if req.Status != "" {
		for _, v := range strings.Split(req.Status, ",") {
			listStatus = append(listStatus, Status(v))
		}
	} else {
		listStatus = []Status{StatusFailed, StatusPending}
	}

	var transactionType []TransactionType
	if req.TransactionType != "" {
		for _, v := range strings.Split(req.TransactionType, ",") {
			transactionType = append(transactionType, TransactionType(v))
		}
	}

	return FilterTransaction{
		Limit:           req.Limit,
		NextCursor:      req.NextCursor,
		PrevCursor:      req.PrevCursor,
		Status:          listStatus,
		TransactionType: transactionType,
	}, nil
}

func (t Transaction) MappingData() Transaction {
	return t
}

func (t Transaction) GetCursorTimestamp() string {
	return fmt.Sprintf("%d", t.Timestamp)
}
