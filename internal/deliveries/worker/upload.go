package worker

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"go-flip-life-style-products/internal/contract"
	"go-flip-life-style-products/internal/models"
	loggerPkg "go-flip-life-style-products/internal/pkg/logger"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	c *contract.Contract
}

func NewUploadWorker(c *contract.Contract) *Worker {
	return &Worker{
		c: c,
	}
}

func (w *Worker) Start(ctx context.Context, size int) {
	loggerPkg.Info(ctx, fmt.Sprintf("Upload Worker started with %d workers", size))
	for i := 0; i < size; i++ {
		go w.run(ctx, i)
	}
}

func (w *Worker) run(ctx context.Context, id int) {
	for job := range w.c.Queue.Channel() {
		loggerPkg.Info(ctx, "worker running upload", logrus.Fields{
			"upload_id": job.UploadID,
			"binary":    id,
		})
		_ = w.store(job)
	}
}

func (w *Worker) store(queue models.UploadQueue) error {
	loggerPkg.Info(queue.Ctx, "processing upload", logrus.Fields{
		"upload_id": queue.UploadID,
	})

	f, err := os.Open(queue.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	header, err := reader.Read()
	if err != nil {
		loggerPkg.Error(queue.Ctx, "failed to read header", logrus.Fields{
			"error": err.Error(),
		})
		return fmt.Errorf("failed to read header: %w", err)
	}

	expectedHeader := []string{
		"timestamp",
		"counterparty",
		"type",
		"amount",
		"status",
		"description",
	}

	if !reflect.DeepEqual(header, expectedHeader) {
		err = fmt.Errorf("invalid CSV header, expected %v, got %v", expectedHeader, header)
		loggerPkg.Error(queue.Ctx, "failed to read header", logrus.Fields{
			"error": err.Error(),
		})
		return err
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil || len(record) != 6 {
			continue
		}

		ts, err := strconv.ParseInt(strings.TrimSpace(record[0]), 10, 64)
		if err != nil {
			continue
		}

		amount, err := decimal.NewFromString(strings.TrimSpace(record[3]))
		if err != nil {
			continue
		}

		err = w.c.Service.Transaction.Store(queue.Ctx, models.Transaction{
			UploadID:     queue.UploadID,
			Timestamp:    ts,
			Counterparty: strings.TrimSpace(record[1]),
			Type:         models.TransactionType(strings.TrimSpace(record[2])),
			Amount:       amount,
			Status:       models.Status(strings.TrimSpace(record[4])),
			Description:  strings.TrimSpace(record[5]),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
