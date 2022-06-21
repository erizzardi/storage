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
// some sensitive values shouldn't have defaults.
// they'd need to be read from k8s secrets
const (
	defaultHTTPPort      = "8081"
	defaultSSLMode       = "disable"
	defaultLogLevel      = "INFO"
	defaultStorageFolder = "./file-storage" // absolute path
	defaultDBDriver      = "postgres"
	defaultDBUsername    = "postgres"  // secret
	defaultDBPassword    = "password"  // secret
	defaultDBDatabase    = "metadata"  // secret
	defaultDBIP          = "localhost" // secret
	defaultDBPort        = "5432"      // secret
	defaultDBTable       = "meta"
)

// global variables, read from environment
var (
	// variables with default
	httpPort          = util.EnvString("STORAGE_HTTP_PORT", defaultHTTPPort)
	sslMode           = util.EnvString("STORAGE_SSL_MODE", defaultSSLMode)
	mainLogLevel      = util.EnvString("STORAGE_MAIN_LOG_LEVEL", defaultLogLevel)
	serviceLogLevel   = util.EnvString("STORAGE_SERVICE_LOG_LEVEL", defaultLogLevel)
	transportLogLevel = util.EnvString("STORAGE_TRANSPORT_LOG_LEVEL", defaultLogLevel)
	endpointsLogLevel = util.EnvString("STORAGE_ENDPOINTS_LOG_LEVEL", defaultLogLevel)
	databaseLogLevel  = util.EnvString("STORAGE_DB_LOG_LEVEL", defaultLogLevel)
	storageFolder     = util.EnvString("STORAGE_FOLDER", defaultStorageFolder)
	dbDriver          = util.EnvString("STORAGE_DB_DRIVER", defaultDBDriver)
	// dbTable           = util.EnvString("STORAGE_DB_TABLE", defaultDBTable)

	// variables without default (secrets)
	dbUser     = os.Getenv("STORAGE_DB_USER")
	dbPassword = os.Getenv("STORAGE_DB_PASSWORD")
	dbDatabase = os.Getenv("STORAGE_DB_DB")
	dbHost     = os.Getenv("STORAGE_DB_HOST")
	dbPort     = os.Getenv("STORAGE_DB_PORT")
)

var (
	// Loggers for every layer.
	// Logging settings can be set individually for each layer
	mainLogger      = util.NewLogger()
	serviceLogger   = util.NewLogger()
	transportLogger = util.NewLogger()
	endpointsLogger = util.NewLogger()
	databaseLogger  = util.NewLogger()
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

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbDatabase, sslMode)
	mainLogger.Debug(connStr)

	switch dbDriver {
	case "postgres":
		db = base.NewSqlDatabase(databaseLogger)
	// TODO - case "mysql":
	// TODO - case "cockroach":
	// TODO - case "cassandra":
	default:
		mainLogger.Error("Unsupported DB type. Supported types: postgres")
		os.Exit(1)
	}
	mainLogger.Info("Connecting to database")
	err := db.Connect(dbDriver, connStr)
	if err != nil {
		mainLogger.Fatal("Error: cannot connect to database: " + err.Error())
		os.Exit(1)
	}
	mainLogger.Info("Database connected")
	defer db.Close() // this fails if db connection is not established

	err = db.Init()
	if err != nil {
		mainLogger.Fatal("Error: cannot initialize database: " + err.Error())
		os.Exit(1)
	}

	//----------------------------------
	// Logging and server initialization
	//----------------------------------
	mainLogger.Debugf("Config variables: %+v\n", config) // TODO

	// All the loggers are passed to the service, so the logging level can be set ar runtime
	var service = storage.NewService(db, serviceLogger, map[string]*util.Logger{
		"main":      mainLogger,
		"transport": transportLogger,
		"endpoints": endpointsLogger,
		"database":  databaseLogger,
	})
	var endpointSet = endpoints.NewEndpointSet(service, config, endpointsLogger)
	var httpHandler = storage.TransportMiddleware{Logger: transportLogger, Next: transport.NewHTTPHandler(endpointSet)}

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
	mainLogger.Warn("Exit: ", g.Run())
}

// init loggers
func init() {
	util.InitLogger(mainLogger, mainLogLevel, logrus.Fields{"level": "main"})
	util.InitLogger(serviceLogger, serviceLogLevel, logrus.Fields{"level": "service"})
	util.InitLogger(transportLogger, transportLogLevel, logrus.Fields{"level": "transport"})
	util.InitLogger(endpointsLogger, endpointsLogLevel, logrus.Fields{"level": "endpoints"})
	util.InitLogger(databaseLogger, databaseLogLevel, logrus.Fields{"level": "database"})
}
