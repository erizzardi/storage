package endpoints

import "io"

type WriteFileRequest struct {
	FileName string
	File     io.Reader
	Err      error
}

type WriteFileResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetFileRequest struct {
	FileName string
}

type GetFileResponse struct {
	Code int    `json:"code"`
	File []byte `json:"-"`
}
