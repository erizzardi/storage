package storage

import (
	"net/http"

	"github.com/erizzardi/storage/util"
)

//------------------
// Transport logging
//------------------
type TransportMiddleware struct {
	Logger *util.Logger
	Next   http.Handler
}

// Middleware for transport layer. It logs every incoming transaction.
// Maybe not needed? Istio/Nginx could do this for free
func (mw TransportMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Logs every incoming request
	mw.Logger.Infof("Incoming request: %s %s", r.Method, r.URL.String())

	// Sets content-type header for every request. If different, it has to be set in the appropriate decodeResponse function
	w.Header().Set("Content-Type", "application/json")

	mw.Next.ServeHTTP(w, r)
}
