package storage

import (
	"net/http"

	"github.com/erizzardi/storage/util"
)

//------------------
// Transport logging
//------------------
type TransportLoggingMiddleware struct {
	Logger *util.Logger
	Next   http.Handler
}

func (mw TransportLoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mw.Logger.Infof("Incoming request: %s\n", r.Method)

	mw.Next.ServeHTTP(w, r)
}
