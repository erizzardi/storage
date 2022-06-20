package base

import (
	"fmt"
	"os"
	"testing"

	"github.com/erizzardi/storage/util"
	"github.com/google/uuid"
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
	//================
	// End of test run
	//================

	if err := db.tearDown(); err != nil {
		testLogger.Error("Cannot drop tables: " + err.Error())
	}
	os.Exit(exitVal)
}

func TestInsertDelete(t *testing.T) {

	var ret util.Row
	uuid := uuid.New().String()
	fileName := "testFile"

	if err := db.InsertMetadata(util.Row{
		Uuid:     uuid,
		FileName: fileName,
	}); err != nil {
		testLogger.Error("Cannot insert row: " + err.Error())
	}
	// check if row was inserted
	statementString := fmt.Sprintf("SELECT * FROM 'meta' WHERE fileName='%s'", fileName)
	rows, err := db.Query(statementString)
	if err != nil {
		testLogger.Error("Cannot query db: " + err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&ret.Uuid, &ret.FileName)
		if err != nil {
			testLogger.Error("Cannot scan row: " + err.Error())
		}
	}
	err = rows.Err()
	if err != nil {
		testLogger.Error("Error: " + err.Error())
	}

	if ret.Uuid != uuid {
		t.Errorf("Uuid not matching:\nSource: %s\nRead:%s", uuid, ret.Uuid)
	}
	if ret.FileName != fileName {
		t.Errorf("fileName not matching:\nSource: %s\nRead:%s", fileName, ret.FileName)
	}

}
