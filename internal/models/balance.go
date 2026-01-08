package models

import "github.com/shopspring/decimal"

type Balance struct {
	UploadID string          `json:"uploadID,omitempty" query:"upload_id"`
	Balance  decimal.Decimal `json:"balance"`
}
