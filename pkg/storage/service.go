package storage

import (
	"context"
	"io"

	"github.com/erizzardi/storage/util"
)

type Service interface {
	//
	//
	// ListFiles list all files paging the request by limit and offset
	ListFiles(ctx context.Context, limit uint, offset uint) ([]util.Row, error)
	//
	//
	// WriteFile writes a file to disk, saving the metadata into the database
	WriteFile(ctx context.Context, file io.Reader, metadata util.Metadata, storageFolder string) (string, error)
	//
	//
	// GetFile gets a file by UUID
	GetFile(ctx context.Context, uuid string, storageFolder string) ([]byte, error)
	//
	//
	// DeleteFile deletes a file by UUID
	DeleteFile(ctx context.Context, uuid string, storageFolder string) error
	//
	//
	// SetLogLevel sets the logging level per layer at runtime
	SetLogLevel(ctx context.Context, layer string, level string) error
}
