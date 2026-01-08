package services_test

import (
	"context"
	"fmt"
	helperPkg "go-flip-life-style-products/internal/pkg/tester"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func fakeMultipartFile() *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: "test.csv",
		Size:     10,
	}
}

func TestUpload(t *testing.T) {
	var (
		unitTestHelper = helperPkg.UnitTestHelper(t)
		ctx            = context.Background()
	)

	type args struct {
		ctx  context.Context
		file *multipart.FileHeader
	}

	tests := []struct {
		name           string
		args           args
		wantErr        bool
		expectedError  error
		expectedOutput string
		doMockService  func(helperPkg.TestHelper, args)
	}{
		{
			name: "success",
			args: args{
				ctx:  ctx,
				file: fakeMultipartFile(),
			},
			wantErr:        false,
			expectedError:  nil,
			expectedOutput: "random-uuid",
			doMockService: func(helper helperPkg.TestHelper, args args) {
				helper.MockFile.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				helper.MockQueue.EXPECT().Enqueue(gomock.Any()).Return(nil)
			},
		},
		{
			name: "error enqueue job",
			args: args{
				ctx:  ctx,
				file: fakeMultipartFile(),
			},
			wantErr:        true,
			expectedError:  fmt.Errorf("got an error"),
			expectedOutput: "",
			doMockService: func(helper helperPkg.TestHelper, args args) {
				helper.MockFile.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
				helper.MockQueue.EXPECT().Enqueue(gomock.Any()).Return(fmt.Errorf("got an error"))
			},
		},
		{
			name: "error save temporary file",
			args: args{
				ctx:  ctx,
				file: fakeMultipartFile(),
			},
			wantErr:        true,
			expectedError:  fmt.Errorf("got an error"),
			expectedOutput: "",
			doMockService: func(helper helperPkg.TestHelper, args args) {
				helper.MockFile.EXPECT().Save(gomock.Any(), gomock.Any()).Return(fmt.Errorf("got an error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.doMockService(unitTestHelper, tt.args)

			res, err := unitTestHelper.StatementsServices.Upload(tt.args.ctx, tt.args.file)
			if (err != nil) == tt.wantErr {
				assert.Equal(t, tt.expectedError, err)
			}

			if !tt.wantErr {
				assert.NoError(t, tt.expectedError, err)
				// cannot assert.equal because the uploadID is random generate
				assert.NotEmpty(t, res)
			}
		})
	}
}
