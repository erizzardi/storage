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
	"github.com/oklog/oklog/pkg/group"
	"github.com/sirupsen/logrus"
)

// defaults
const (
	defaultHTTPPort = "8081"
	defaultLogLevel = "INFO"
)

// global variables, read from environment
var (
	httpPort          = util.EnvString("STORAGE_HTTP_PORT", defaultHTTPPort)
	serviceLogLevel   = util.EnvString("STORAGE_SERVICE_LOG_LEVEL", defaultLogLevel)
	transportLogLevel = util.EnvString("STORAGE_TRANSPORT_LOG_LEVEL", defaultLogLevel)
)

// Loggers for application and transport layer
// Logging settings can be set for each layer individually
var (
	serviceLogger   = logrus.New()
	transportLogger = logrus.New()
)

func main() {
	var httpAddr = net.JoinHostPort("localhost", httpPort)

	serviceLogger.Info("Service started. Listening from port " + httpPort)

	var service = storage.NewService()
	service = storage.ServiceLoggingMiddleware{Logger: serviceLogger, Next: service}

	var endpointSet = endpoints.NewEndpointSet(service)

	var httpHandler = transport.NewHTTPHandler(endpointSet)
	httpHandler = storage.TransportLoggingMiddleware{Logger: transportLogger, Next: httpHandler}

	var g group.Group
	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			os.Exit(1)
		}
		g.Add(func() error {
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
	serviceLogger.Info("Exit: ", g.Run())
}

// init logrus
func init() {
	serviceLogger.SetFormatter(&logrus.JSONFormatter{})
	transportLogger.SetFormatter(&logrus.JSONFormatter{})

	// loglevel for service logger
	switch serviceLogLevel {
	case "DEBUG":
		serviceLogger.SetLevel(logrus.DebugLevel)
	case "INFO":
		serviceLogger.SetLevel(logrus.InfoLevel)
	case "WARN":
		serviceLogger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		serviceLogger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		serviceLogger.SetLevel(logrus.FatalLevel)
	}

	// loglevel for transport logger
	switch transportLogLevel {
	case "DEBUG":
		transportLogger.SetLevel(logrus.DebugLevel)
	case "INFO":
		transportLogger.SetLevel(logrus.InfoLevel)
	case "WARN":
		transportLogger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		transportLogger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		transportLogger.SetLevel(logrus.FatalLevel)
	}

}
