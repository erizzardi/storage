package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/erizzardi/storage/pkg/storage/endpoints"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	m := http.NewServeMux()

	m.Handle("/writefile", httptransport.NewServer(
		ep.WriteFileEndpoint,
		decodeHTTPWriteFileRequest,
		encodeResponse,
	))

	m.Handle("/getfile", httptransport.NewServer(
		ep.GetFileEndpoint,
		decodeHTTPGetFileRequest,
		encodeGetFileResponse,
	))

	return m
}

func decodeHTTPWriteFileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	file, _, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		return nil, err
	}
	fileName := r.FormValue("file-name")
	return endpoints.WriteFileRequest{
		FileName: fileName,
		File:     file,
	}, nil
}

func decodeHTTPGetFileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	fileName := r.FormValue("file-name")
	return endpoints.GetFileRequest{
		FileName: fileName,
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeGetFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoints.GetFileResponse)

	w.Write(res.File)
	w.Header().Set("Content-Type", "application/octet-stream")

	return json.NewEncoder(w).Encode(response)
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
