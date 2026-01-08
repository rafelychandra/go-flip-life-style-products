package file

import (
	"io"
	"mime/multipart"
	"os"
)

type (
	File interface {
		Save(file *multipart.FileHeader, dst string) error
	}

	file struct{}
)

func NewFile() File {
	return &file{}
}

func (f *file) Save(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
