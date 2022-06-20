package base

import (
	"os"
	"testing"

	"github.com/erizzardi/storage/util"
)

// Unit tests for the postgres implementation of the db interface.

var (
	testDbConnStr = os.Getenv("TEST_DB_CONN_STR")
	testLogger    = util.NewLogger()
	db            = NewSqlDatabase(testLogger) // change

	driver = "postgres"
)

func TestMain(m *testing.M) {

	if err := db.Connect(driver, testDbConnStr); err != nil {
		testLogger.Error("Cannot connect to database: " + err.Error())
	}
	if err := db.Init(); err != nil {
		testLogger.Error("Cannot init database: " + err.Error())
	}
	//==================================
	// Function that actually runs tests
	//==================================
	exitVal := m.Run()
	if err := db.tearDown(); err != nil {
		testLogger.Error("Cannot drop tables: " + err.Error())
	}
	os.Exit(exitVal)
}

func TestInsertDelete(t *testing.T) {

}
