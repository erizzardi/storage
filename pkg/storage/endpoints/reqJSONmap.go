package endpoints

import (
	"io"
	"net/http"

	"github.com/erizzardi/storage/util"
)

type Request interface{}

type Response interface{}

//=========
// Requests
//=========

type NotFoundRequest struct {
	Endpoint string
}

type WriteFileRequest struct {
	File     io.Reader
	Metadata util.Metadata
	Headers  http.Header
}

type GetFileRequest struct {
	Uuid    string
	Headers http.Header
}

type DeleteFileRequest struct {
	Uuid    string
	Headers http.Header
}

type AddBucketRequest struct {
	Name       string `json:"name"`
	Versioning bool   `json:"versioning"`
	Headers    http.Header
	// TODO - LifecyclePolicy util.LifecyclePolicy
}

type LogLevelRequest struct {
	Layer   string `json:"layer"`
	Level   string `json:"level"`
	Headers http.Header
}

type ListFilesRequest struct {
	Limit   uint `json:"limit"`
	Offset  uint `json:"offset"`
	Headers http.Header
}

//==========
// Responses
//==========

type HealtzResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type NotFoundResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WriteFileResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Uuid    string `json:"uuid,omitempty"`
}

type GetFileResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	File    []byte `json:"-"`
}

type DeleteFileResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AddBucketResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LogLevelResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ListFilesResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Files   []util.Row `json:"files"`
}
