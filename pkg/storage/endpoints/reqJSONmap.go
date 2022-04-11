package endpoints

import (
	"io"

	"github.com/erizzardi/storage/util"
)

type HealtzResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WriteFileRequest struct {
	File     io.Reader
	Metadata util.Metadata
}

type WriteFileResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Uuid    string `json:"uuid"`
}

type GetFileRequest struct {
	Uuid string
}

type GetFileResponse struct {
	Code int    `json:"code"`
	File []byte `json:"-"`
}

type DeleteFileRequest struct {
	Uuid string
}

type DeleteFileResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AddBucketRequest struct {
	Name       string `json:"name"`
	Versioning bool   `json:"versioning"`
	// TODO - LifecyclePolicy util.LifecyclePolicy
}

type AddBucketResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
