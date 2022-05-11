package base

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/erizzardi/storage/util"
)

var db, mock, _ = sqlmock.New()

var table = "table"

var testDB = NewSqlDatabase(db, util.NewLogger(), table)

// ????
func TestExecSuccess(t *testing.T) {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO "+table).WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	testDB.Exec("INSERT INTO " + table + " VALUES(1,2)")
}
