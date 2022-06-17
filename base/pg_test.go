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
	db            = NewSqlDatabase(testLogger, "meta") // change

	driver = "postgres"
)

func TestConnect(t *testing.T) {

	err := db.Connect(driver, testDbConnStr)
	if err != nil {
		t.Error("Cannot connect to database: " + err.Error())
	}

}
