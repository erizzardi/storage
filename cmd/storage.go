package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/erizzardi/storage/pkg/storage"
	"github.com/erizzardi/storage/pkg/storage/endpoints"
	"github.com/erizzardi/storage/pkg/storage/transport"
	"github.com/erizzardi/storage/util"
	"github.com/go-kit/log"
	"github.com/oklog/oklog/pkg/group"
)

const (
	defaultHTTPPort = "8081"
)

func main() {
	var httpAddr = net.JoinHostPort("localhost", util.EnvString("STORAGE_HTTP_PORT", defaultHTTPPort))

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	var (
		service     = storage.NewService()
		endpointSet = endpoints.NewEndpointSet(service)
		httpHandler = transport.NewHTTPHandler(endpointSet)
	)

	fmt.Printf("%+v\n", endpointSet.GetFileEndpoint)

	var g group.Group
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	logger.Log("exit", g.Run())
}
