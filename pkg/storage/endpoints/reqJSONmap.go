package endpoints

import (
	"io"

	"github.com/erizzardi/storage/util"
)

type Request interface{}

type Response interface{}

//=========
// Requests
//=========

type WriteFileRequest struct {
	File     io.Reader
	Metadata util.Metadata
}

type GetFileRequest struct {
	Uuid string
}

type DeleteFileRequest struct {
	Uuid string
}

type AddBucketRequest struct {
	Name       string `json:"name"`
	Versioning bool   `json:"versioning"`
	// TODO - LifecyclePolicy util.LifecyclePolicy
}

type LogLevelRequest struct {
	Layer string `json="layer"`
	Level string `json="level"`
}

//==========
// Responses
//==========

type HealtzResponse struct {
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
	Message string `json="message"`
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
