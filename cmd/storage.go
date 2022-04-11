package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/erizzardi/storage/base"
	"github.com/erizzardi/storage/pkg/storage"
	"github.com/erizzardi/storage/pkg/storage/endpoints"
	"github.com/erizzardi/storage/pkg/storage/transport"
	"github.com/erizzardi/storage/util"
	"github.com/oklog/oklog/pkg/group"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

// defaults
// TODO - change with better names!!
const (
	defaultHTTPPort      = "8081"
	defaultLogLevel      = "INFO"
	defaultStorageFolder = "./file-storage" // absolute path
	defaultDBDriver      = "postgres"
	defaultDBUsername    = "postgres"
	defaultDBPassword    = "postgres"
	defaultDBDatabase    = "storage-metadata"
	defaultDBIP          = "localhost"
	defaultDBPort        = "5432"
	defaultDBTable       = "meta"
)

// global variables, read from environment
var (
	httpPort          = util.EnvString("STORAGE_HTTP_PORT", defaultHTTPPort)
	mainLogLevel      = util.EnvString("STORAGE_MAIN_LOG_LEVEL", defaultLogLevel)
	serviceLogLevel   = util.EnvString("STORAGE_SERVICE_LOG_LEVEL", defaultLogLevel)
	transportLogLevel = util.EnvString("STORAGE_TRANSPORT_LOG_LEVEL", defaultLogLevel)
	databaseLogLevel  = util.EnvString("STORAGE_DB_LOG_LEVEL", defaultLogLevel)
	storageFolder     = util.EnvString("STORAGE_FOLDER", defaultStorageFolder)
	dbDriver          = util.EnvString("STORAGE_DB_DRIVER", defaultDBDriver)
	dbUser            = util.EnvString("STORAGE_DB_USER", defaultDBUsername)
	dbPassword        = util.EnvString("STORAGE_DB_PASSWORD", defaultDBPassword)
	dbDatabase        = util.EnvString("STORAGE_DB_DB", defaultDBDatabase)
	dbHost            = util.EnvString("STORAGE_DB_HOST", defaultDBIP)
	dbPort            = util.EnvString("STORAGE_DB_PORT", defaultDBPort)
	dbTable           = util.EnvString("STORAGE_DB_TABLE", defaultDBTable)
)

// Loggers for every layer.
// Logging settings can be set individually for each layer
var (
	mainLogger      = logrus.New()
	serviceLogger   = logrus.New()
	transportLogger = logrus.New()
	databaseLogger  = logrus.New()
)

func main() {
	//--------------
	// Set up config
	//--------------
	var config = util.SetConfig(storageFolder)
	// Listening HTTP address
	var httpAddr = net.JoinHostPort("localhost", httpPort)

	//---------------------------------
	// DB connection and initialization
	//---------------------------------
	var db base.DB

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbDatabase)
	mainLogger.Debug(connStr)

	switch dbDriver {
	case "postgres":
		db = base.NewSqlDatabase(databaseLogger, dbTable)
	// TODO - case "mysql":
	// TODO - case "cassandra":
	default:
		mainLogger.Error("Unsupported DB type. Supported types: postgres")
		os.Exit(1)
	}
	err := db.Connect(dbDriver, connStr)
	if err != nil {
		mainLogger.Error("Error: cannot connect to database: " + err.Error())
	}
	defer db.Close()

	err = db.Init()
	if err != nil {
		mainLogger.Error("Error: cannot connect to database: " + err.Error())
	}

	//----------------------------------
	// Logging and server initialization
	//----------------------------------
	mainLogger.Debugf("Config variables: %+v\n", config) // TODO

	// var service = storage.ServiceLoggingMiddleware{Logger: serviceLogger, Next: storage.NewService()}
	var service = storage.NewService(db, serviceLogger)
	var endpointSet = endpoints.NewEndpointSet(service, config)
	var httpHandler = storage.TransportLoggingMiddleware{Logger: transportLogger, Next: transport.NewHTTPHandler(endpointSet)}

	mainLogger.Info("Service initialization complete. Listening on port " + httpPort)

	//-----------------------------
	// Run HTTP listener and server
	//-----------------------------
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
	// {
	/*
		TODO - Implement object lifecycle
	*/
	// }
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
	mainLogger.Info("Exit: ", g.Run())
}

// init logrus
func init() {
	// mainLogger.SetFormatter(&logrus.JSONFormatter{})
	// serviceLogger.SetFormatter(&logrus.JSONFormatter{})
	// transportLogger.SetFormatter(&logrus.JSONFormatter{})
	// databaseLogger.SetFormatter(&logrus.JSONFormatter{})
	mainLogger.WithFields(logrus.Fields{
		"Layer": "main",
	})
	serviceLogger.WithFields(logrus.Fields{
		"Layer": "service",
	})
	transportLogger.WithFields(logrus.Fields{
		"Layer": "transport",
	})
	databaseLogger.WithFields(logrus.Fields{
		"Layer": "database",
	})

	// loglevel for main logger
	switch mainLogLevel {
	case "DEBUG":
		mainLogger.SetLevel(logrus.DebugLevel)
	case "INFO":
		mainLogger.SetLevel(logrus.InfoLevel)
	case "WARN":
		mainLogger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		mainLogger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		serviceLogger.SetLevel(logrus.FatalLevel)
	}

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

	// loglevel for database logger
	switch databaseLogLevel {
	case "DEBUG":
		databaseLogger.SetLevel(logrus.DebugLevel)
	case "INFO":
		databaseLogger.SetLevel(logrus.InfoLevel)
	case "WARN":
		databaseLogger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		databaseLogger.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		databaseLogger.SetLevel(logrus.FatalLevel)
	}
}
