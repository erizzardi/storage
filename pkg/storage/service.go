package storage

import (
	"context"
	"io"

	"github.com/erizzardi/storage/util"
)

type Service interface {
	WriteFile(ctx context.Context, file io.Reader, metadata util.Metadata, storageFolder string) (string, error)
	GetFile(ctx context.Context, uuid string, storageFolder string) ([]byte, error)
	DeleteFile(ctx context.Context, uuid string, storageFolder string) error
}
