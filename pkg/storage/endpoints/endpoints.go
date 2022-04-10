package endpoints

import (
	"context"
	"os"

	"github.com/erizzardi/storage/pkg/storage"
	"github.com/erizzardi/storage/util"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

type Set struct {
	HealtzEndpoint    endpoint.Endpoint
	WriteFileEndpoint endpoint.Endpoint
	GetFileEndpoint   endpoint.Endpoint
}

//----------------------------------------
// No logging here!! Need to investigate!!
//----------------------------------------
func NewEndpointSet(svc storage.Service, config *util.Config) Set {
	return Set{
		HealtzEndpoint:    MakeHealtzEndpoint(),
		WriteFileEndpoint: MakeWriteFileEndpoint(svc, config.StorageFolder),
		GetFileEndpoint:   MakeGetFileEndpoint(svc, config.StorageFolder),
	}
}

func MakeHealtzEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return HealtzResponse{Message: "Alive!", Code: 200}, nil
	}
}

func MakeWriteFileEndpoint(svc storage.Service, storageFolder string) endpoint.Endpoint {
	//possibly cluster all config variables in one struct and pass that to the WriteFile method
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(WriteFileRequest)
		err := svc.WriteFile(ctx, req.File, req.FileName, storageFolder)
		if err != nil {
			er := err.(*util.ResponseError)
			return WriteFileResponse{er.StatusCode, err.Error()}, nil
		}
		return WriteFileResponse{201, "File created"}, nil
	}
}

func MakeGetFileEndpoint(svc storage.Service, storageFolder string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetFileRequest)
		file, err := svc.GetFile(ctx, req.FileName, storageFolder)
		if err != nil {
			er := err.(*util.ResponseError)
			return GetFileResponse{er.StatusCode, nil}, err
		}
		return GetFileResponse{200, file}, nil
	}
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
