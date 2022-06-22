package base

import (
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

//
// This test inserts, retrieves and deletes a row of metadata.
// Pass if no errors.
func TestInsertRetrieveDelete(t *testing.T) {

	var ret util.Row
	// count := 0

	uuid := uuid.New().String()
	fileName := "testFile"

	//----------------------------
	// insert (uuid, testFile) row
	//----------------------------
	if err := db.InsertMetadata(util.Row{
		Uuid:     uuid,
		FileName: fileName,
	}); err != nil {
		t.Error("Cannot insert row: " + err.Error())
	}

	//----------------------------------------
	// check if row exists - retrieve metadata
	//----------------------------------------
	ret, err := db.RetrieveMetadata("uuid", uuid)
	if err != nil {
		t.Error("Cannot insert row: " + err.Error())
	}

	//-----------
	// assertions
	//-----------
	if ret.Uuid != uuid {
		t.Errorf("Uuid not matching:\nSource: %s\nRead:%s", uuid, ret.Uuid)
	}
	if ret.FileName != fileName {
		t.Errorf("fileName not matching:\nSource: %s\nRead:%s", fileName, ret.FileName)
	}

	//-----------
	// delete row
	//-----------
	// statementString := fmt.Sprintf("DELETE FROM meta WHERE filename='%s';", fileName)
	// if _, err := db.Exec(statementString); err != nil {
	// 	t.Error("Error deleting the row: " + err.Error())
	// }
	if err = db.DeleteMetadata("uuid", uuid); err != nil {
		t.Error(err.Error())
	}

	//-------------------------------
	// check if the row doesn't exist
	//-------------------------------
	ret, err = db.RetrieveMetadata("uuid", uuid)
	if err != nil {
		t.Error("Cannot insert row: " + err.Error())
	}

	if (ret != util.Row{}) {
		t.Error("Found Row! Data not deleted correctly.")
	}
}

//
// This test deletes a non-existing row of metadata.
// Pass if errors
func TestDeleteNonexisting(t *testing.T) {

	// Every uuid is unique by construction
	// thus creating a new one is sufficient
	uuid := uuid.New().String()

	//-----------
	// delete row
	//-----------
	if err := db.DeleteMetadata("uuid", uuid); err == nil {
		t.Error("Should return error")
	} else if util.ErrorIs(err, util.BadRequestError{}) == false {
		t.Errorf("Error type should be %T", util.BadRequestError{})
	}
}
