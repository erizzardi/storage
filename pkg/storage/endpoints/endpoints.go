package endpoints

import (
	"context"
	"errors"

	"github.com/erizzardi/storage/pkg/storage"
	"github.com/erizzardi/storage/util"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	HealtzEndpoint           endpoint.Endpoint
	NotFoundEndpoint         endpoint.Endpoint
	MethodNotAllowedEndpoint endpoint.Endpoint
	WriteFileEndpoint        endpoint.Endpoint
	GetFileEndpoint          endpoint.Endpoint
	DeleteFileEndpoint       endpoint.Endpoint
	AddBucketEndpoint        endpoint.Endpoint
	LogLevelEndpoint         endpoint.Endpoint
	ListFilesEndpoint        endpoint.Endpoint
}

func NewEndpointSet(svc storage.Service, config *util.Config, logger *util.Logger) Set {
	return Set{
		HealtzEndpoint:           MakeHealtzEndpoint(logger),
		NotFoundEndpoint:         MakeNotFoundEndpoint(logger),
		MethodNotAllowedEndpoint: MakeMethodNotAllowedEndpoint(logger),
		WriteFileEndpoint:        MakeWriteFileEndpoint(svc, config.StorageFolder, logger),
		GetFileEndpoint:          MakeGetFileEndpoint(svc, config.StorageFolder, logger),
		DeleteFileEndpoint:       MakeDeleteFileEndpoint(svc, config.StorageFolder, logger),
		AddBucketEndpoint:        MakeAddBucketEndpoint(svc, config.StorageFolder, logger),
		LogLevelEndpoint:         MakeLogLevelEndpoint(svc, config.StorageFolder, logger),
		ListFilesEndpoint:        MakeListFilesEndpoint(svc, config.StorageFolder, logger),
	}
}

//=================================================================================
// Endpoint layer
//---------------------------------------------------------------------------------
// NEVER return errors via endpoint.Endpoint. It fucks up all the response encoding
//=================================================================================

func MakeHealtzEndpoint(logger *util.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return HealtzResponse{Message: "Alive!", Code: 200}, nil
	}
}

func MakeNotFoundEndpoint(logger *util.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		endpoint := request.(NotFoundRequest).Endpoint
		logger.Error("Requested endpoint not found: " + endpoint)
		return HealtzResponse{Message: "Not found: " + endpoint, Code: 404}, nil
	}
}

func MakeMethodNotAllowedEndpoint(logger *util.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		method := request.(MethodNotAllowedRequest).Method
		logger.Error("Requested method not allowed: " + method)
		return HealtzResponse{Message: "Method not allowed: " + method, Code: 415}, nil
	}
}

func MakeListFilesEndpoint(svc storage.Service, storageFolder string, logger *util.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListFilesRequest)
		if req.Err != nil {
			return WriteFileResponse{400, "Could not read body: " + req.Err.Error(), ""}, nil
		}
		files, err := svc.ListFiles(ctx, req.Limit, req.Offset)
		if errors.Is(err, util.InternalServerError{}) {
			return ListFilesResponse{Code: 500, Message: err.Error(), Files: nil}, nil
		}
		return ListFilesResponse{Code: 200, Message: "Ok", Files: files}, nil
	}
}

func MakeWriteFileEndpoint(svc storage.Service, storageFolder string, logger *util.Logger) endpoint.Endpoint {
	//possibly cluster all config variables in one struct and pass that to the WriteFile method
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(WriteFileRequest)
		uuid, err := svc.WriteFile(ctx, req.File, req.Metadata, storageFolder)
		if err != nil {
			er := err.(*util.ResponseError) // TODO - is it necessary? investigate a more elegant solution
			return WriteFileResponse{er.StatusCode, err.Error(), ""}, nil
		}
		return WriteFileResponse{201, "File created", uuid}, nil
	}
}

func MakeGetFileEndpoint(svc storage.Service, storageFolder string, logger *util.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetFileRequest)
		file, err := svc.GetFile(ctx, req.Uuid, storageFolder)
		if err != nil {
			er := err.(*util.ResponseError) // TODO - as above
			return GetFileResponse{er.StatusCode, err.Error(), nil}, nil
		}
		return GetFileResponse{200, "File retrieved", file}, nil
	}
}

func MakeDeleteFileEndpoint(svc storage.Service, storageFolder string, logger *util.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteFileRequest)
		err := svc.DeleteFile(ctx, req.Uuid, storageFolder)
		if err != nil {
			er := err.(*util.ResponseError) // TODO - as above
			return DeleteFileResponse{er.StatusCode, err.Error()}, nil
		}
		return DeleteFileResponse{200, "File deleted"}, nil
	}
}

func MakeAddBucketEndpoint(svc storage.Service, storageFolder string, logger *util.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// req := request.(AddBucketRequest)
		return AddBucketResponse{200, "Bucket created"}, nil
	}
}

func MakeLogLevelEndpoint(svc storage.Service, storageFolder string, logger *util.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LogLevelRequest)
		err := svc.SetLogLevel(ctx, req.Layer, req.Level)
		if err != nil {
			er := err.(*util.ResponseError)
			return LogLevelResponse{er.StatusCode, err.Error()}, nil
		}
		return LogLevelResponse{200, "Logging level for layer " + req.Layer + " changed to " + req.Level}, nil
	}
}
