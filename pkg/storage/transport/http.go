package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/erizzardi/storage/pkg/storage/endpoints"
	"github.com/erizzardi/storage/util"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	r := mux.NewRouter()

	r.NotFoundHandler = httptransport.NewServer(
		ep.NotFoundEndpoint,
		decodeHTTPNotFoundRequest,
		encodeNotFoundResponse,
	)

	r.Methods("GET").Path("/healtz").Handler(httptransport.NewServer(
		ep.HealtzEndpoint,
		decodeHTTPHealtzRequest,
		encodeHealthzResponse,
	))

	r.Methods("GET").Path("/files").Handler(httptransport.NewServer(
		ep.ListFilesEndpoint,
		decodeHTTPListFilesRequest,
		encodeListFilesResponse,
	))

	r.Methods("POST").Path("/files").Handler(httptransport.NewServer(
		ep.WriteFileEndpoint,
		decodeHTTPWriteFileRequest,
		encodeWriteFileResponse,
	))

	r.Methods("GET").Path("/files/{id}").Handler(httptransport.NewServer(
		ep.GetFileEndpoint,
		decodeHTTPGetFileRequest,
		encodeGetFileResponse,
	))

	r.Methods("DELETE").Path("/files/{id}").Handler(httptransport.NewServer(
		ep.DeleteFileEndpoint,
		decodeHTTPDeleteFileRequest,
		encodeDeleteFileResponse,
	))

	r.Methods("PUT").Path("/buckets").Handler(httptransport.NewServer(
		ep.AddBucketEndpoint,
		decodeHTTPAddBucketRequest,
		encodeAddBucketResponse,
	))

	r.Methods("POST").Path("/config/loglevel").Handler(httptransport.NewServer(
		ep.LogLevelEndpoint,
		decodeHTTPLogLevelRequest,
		encodeLogLevelResponse,
	))

	return r
}

//=================
// Request Decoders
//=================
func decodeHTTPHealtzRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeHTTPNotFoundRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoints.NotFoundRequest{Endpoint: r.URL.String()}, nil
}

func decodeHTTPListFilesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := &endpoints.ListFilesRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	return *req, err
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

func decodeHTTPAddBucketRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := &endpoints.AddBucketRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	return *req, err
}

func decodeHTTPLogLevelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := &endpoints.LogLevelRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	return *req, err
}

//==================
// Response Encoders
//==================
func encodeHealthzResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	return json.NewEncoder(w).Encode(response)
}

func encodeListFilesResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
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
	w.Header().Set("Content-Type", "image/jpg")
	return nil
}

func encodeDeleteFileResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoints.DeleteFileResponse)
	w.WriteHeader(res.Code)
	return json.NewEncoder(w).Encode(response)
}

func encodeAddBucketResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoints.AddBucketResponse)
	w.WriteHeader(res.Code)
	return json.NewEncoder(w).Encode(response)
}

func encodeLogLevelResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(endpoints.LogLevelResponse)
	w.WriteHeader(res.Code)
	return json.NewEncoder(w).Encode(response)
}

func encodeNotFoundResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
