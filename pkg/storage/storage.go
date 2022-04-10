package storage

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type storageService struct{}

func NewService() Service { return &storageService{} }

func (*storageService) WriteFile(ctx context.Context, file io.Reader, fileName string) error {
	fileName = filepath.Join("./", fileName)
	newFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)

	return err
}

func (*storageService) GetFile(ctx context.Context, fileName string) ([]byte, error) {
	fileName = filepath.Join("./", fileName)
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}
