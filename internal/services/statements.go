package services

import (
	"context"
	"go-flip-life-style-products/internal/config"
	"go-flip-life-style-products/internal/models"
	filePkg "go-flip-life-style-products/internal/pkg/file"
	queuePkg "go-flip-life-style-products/internal/pkg/queue"
	uuidPkg "go-flip-life-style-products/internal/pkg/uuid"
	"mime/multipart"
	"os"
	"path/filepath"
)

type (
	Statements interface {
		Upload(ctx context.Context, file *multipart.FileHeader) (string, error)
	}

	statements struct {
		cfg   *config.Configuration
		file  filePkg.File
		queue queuePkg.Queue
	}
)

func NewStatements(cfg *config.Configuration, file filePkg.File, queue queuePkg.Queue) Statements {
	return &statements{
		cfg:   cfg,
		file:  file,
		queue: queue,
	}
}

func (s *statements) Upload(ctx context.Context, file *multipart.FileHeader) (string, error) {
	var (
		uploadID = uuidPkg.UUID()
		tmpDir   = os.TempDir()
		filePath = filepath.Join(tmpDir, uploadID+".csv")
	)

	err := s.file.Save(file, filePath)
	if err != nil {
		return "", err
	}

	err = s.queue.Enqueue(models.UploadQueue{
		Ctx:      ctx,
		UploadID: uploadID,
		FilePath: filePath,
	})
	if err != nil {
		return "", err
	}

	return uploadID, nil
}
