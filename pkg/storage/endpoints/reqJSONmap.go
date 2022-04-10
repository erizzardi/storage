package endpoints

import "io"

type HealtzResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WriteFileRequest struct {
	FileName string
	File     io.Reader
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
