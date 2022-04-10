package endpoints

import (
	"context"
	"os"

	"github.com/erizzardi/storage/pkg/storage"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	WriteFileEndpoint endpoint.Endpoint
	GetFileEndpoint   endpoint.Endpoint
}

func NewEndpointSet(svc storage.Service) Set {
	return Set{
		WriteFileEndpoint: MakeWriteFileEndpoint(svc),
		GetFileEndpoint:   MakeGetFileEndpoint(svc),
	}
}

func MakeWriteFileEndpoint(svc storage.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(WriteFileRequest)
		if req.Err != nil {
			return WriteFileResponse{500, "Internal server error"}, req.Err
		}
		err := svc.WriteFile(ctx, req.File, req.FileName)
		if err != nil {
			return WriteFileResponse{500, "Internal server error"}, err
		}
		return WriteFileResponse{200, "Ok!"}, nil
	}
}

func MakeGetFileEndpoint(svc storage.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetFileRequest)
		file, err := svc.GetFile(ctx, req.FileName)
		if err != nil {
			return GetFileResponse{500, nil}, err
		}
		return GetFileResponse{200, file}, nil
	}
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
