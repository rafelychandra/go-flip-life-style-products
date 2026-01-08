package models

import "context"

type UploadQueue struct {
	Ctx      context.Context
	UploadID string
	FilePath string
}
