package storage

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

//--------------------
// Application logging
//--------------------
type ServiceLoggingMiddleware struct {
	Logger *logrus.Logger
	Next   Service
}

func (mw ServiceLoggingMiddleware) WriteFile(ctx context.Context, file io.Reader, fileName string) error {
	defer func(begin time.Time) {
		mw.Logger.WithFields(logrus.Fields{
			"fileName": fileName,
		}).Info("WriteFile method invoked")
	}(time.Now())

	return mw.Next.WriteFile(ctx, file, fileName)
}

func (mw ServiceLoggingMiddleware) GetFile(ctx context.Context, fileName string) ([]byte, error) {
	defer func(begin time.Time) {
		mw.Logger.WithFields(logrus.Fields{
			"fileName": fileName,
		}).Info("GetFile method invoked")
	}(time.Now())

	return mw.Next.GetFile(ctx, fileName)
}

//------------------
// Transport logging
//------------------

type TransportLoggingMiddleware struct {
	Logger *logrus.Logger
	Next   http.Handler
}

func (mw TransportLoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	mw.Logger.WithFields(logrus.Fields{
		"Request": r.Method + " " + r.URL.Path,
	}).Info("Incoming HTTP request")

	mw.Next.ServeHTTP(w, r)
}
