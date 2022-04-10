package storage

import (
	"context"
	"io"
)

type Service interface {
	WriteFile(ctx context.Context, file io.Reader, storageFolder string) (string, error)
	GetFile(ctx context.Context, uuid string, storageFolder string) ([]byte, error)
	DeleteFile(ctx context.Context, uuid string, storageFolder string) error
}
