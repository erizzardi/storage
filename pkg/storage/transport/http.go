package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/erizzardi/storage/pkg/storage/endpoints"
	"github.com/erizzardi/storage/util"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	r := mux.NewRouter()

	r.Methods("GET").Path("/healtz").Handler(httptransport.NewServer(
		ep.HealtzEndpoint,
		decodeHTTPHealtzRequest,
		encodeHealthResponse,
	))

	r.Methods("POST").Path("/file").Handler(httptransport.NewServer(
		ep.WriteFileEndpoint,
		decodeHTTPWriteFileRequest,
		encodeWriteFileResponse,
	))

	r.Methods("GET").Path("/file/{id}").Handler(httptransport.NewServer(
		ep.GetFileEndpoint,
		decodeHTTPGetFileRequest,
		encodeGetFileResponse,
	))

	r.Methods("DELETE").Path("/file/{id}").Handler(httptransport.NewServer(
		ep.DeleteFileEndpoint,
		decodeHTTPDeleteFileRequest,
		encodeDeleteFileResponse,
	))

	return r
}

func decodeHTTPHealtzRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeHTTPWriteFileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()

	file, multipartHeader, err := r.FormFile("file")
	if err != nil && err != http.ErrMissingFile {
		return nil, err
	}
	return endpoints.WriteFileRequest{
		File: file,
		Metadata: util.Metadata{
			Name: multipartHeader.Filename,
			Size: multipartHeader.Size,
		},
	}, nil
}

func decodeHTTPGetFileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	uuid := vars["id"]

	return endpoints.GetFileRequest{
		Uuid: uuid,
	}, nil
}

func decodeHTTPDeleteFileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	uuid := vars["id"]

	return endpoints.DeleteFileRequest{
		Uuid: uuid,
	}, nil
}

func encodeHealthResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeWriteFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoints.WriteFileResponse)
	w.WriteHeader(res.Code)
	return json.NewEncoder(w).Encode(response)
}

func encodeGetFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoints.GetFileResponse)

	w.Write(res.File)
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(response)
}

func encodeDeleteFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoints.DeleteFileResponse)

	w.WriteHeader(res.Code)
	return json.NewEncoder(w).Encode(response)
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
