package endpoints

import (
	"context"

	"github.com/erizzardi/storage/pkg/storage"
	"github.com/erizzardi/storage/util"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	HealtzEndpoint     endpoint.Endpoint
	WriteFileEndpoint  endpoint.Endpoint
	GetFileEndpoint    endpoint.Endpoint
	DeleteFileEndpoint endpoint.Endpoint
}

//-----------------------------------------------
// TODO - No logging here!! Need to investigate!!
// 		  Probably I need another middleware
//-----------------------------------------------
func NewEndpointSet(svc storage.Service, config *util.Config) Set {
	return Set{
		HealtzEndpoint:     MakeHealtzEndpoint(),
		WriteFileEndpoint:  MakeWriteFileEndpoint(svc, config.StorageFolder),
		GetFileEndpoint:    MakeGetFileEndpoint(svc, config.StorageFolder),
		DeleteFileEndpoint: MakeDeleteFileEndpoint(svc, config.StorageFolder),
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
		uuid, err := svc.WriteFile(ctx, req.File, req.Metadata, storageFolder)
		if err != nil {
			er := err.(*util.ResponseError) // TODO - is it necessary? investigate a more elegant solution
			return WriteFileResponse{er.StatusCode, err.Error(), ""}, nil
		}
		return WriteFileResponse{201, "File created", uuid}, nil
	}
}

func MakeGetFileEndpoint(svc storage.Service, storageFolder string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetFileRequest)
		file, err := svc.GetFile(ctx, req.Uuid, storageFolder)
		if err != nil {
			er := err.(*util.ResponseError) // TODO - as above
			return GetFileResponse{er.StatusCode, nil}, nil
		}
		return GetFileResponse{200, file}, nil
	}
}

func MakeDeleteFileEndpoint(svc storage.Service, storageFolder string) endpoint.Endpoint {
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
