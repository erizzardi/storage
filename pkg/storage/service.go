package storage

import (
	"context"
	"io"
)

type Service interface {
	WriteFile(ctx context.Context, file io.Reader, fileName string) error
	GetFile(ctx context.Context, fileName string) ([]byte, error)
}
