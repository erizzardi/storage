package endpoints

import "io"

type HealtzResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type WriteFileRequest struct {
	FileName string
	File     io.Reader
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
