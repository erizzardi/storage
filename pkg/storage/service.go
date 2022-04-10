package storage

import (
	"context"
	"io"
)

type Service interface {
	WriteFile(ctx context.Context, file io.Reader, fileName string, storageFolder string) error
	GetFile(ctx context.Context, fileName string, storageFolder string) ([]byte, error)
	DeleteFile(ctx context.Context, fileName string, storageFolder string) error
}
