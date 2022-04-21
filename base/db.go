package base

import (
	"database/sql"

	"github.com/erizzardi/storage/util"
)

// Interface for database connection and operations.
// Implementation for sql database: sql.go
type DB interface {
	//
	//
	// Enstablishes connection to the database
	Connect(driver string, dsl string) error
	//
	//
	// Creates the table in the database
	Init() error
	//
	//
	// Prepares and executes 'statement' operation.
	// Use with INSERT, CREATE, DELETE statements
	Exec(statement string, params ...any) (sql.Result, error)
	//
	//
	// Prepares and executes 'statement' operation.
	// Use with SELECT statements
	Query(statement string, params ...any) (*sql.Rows, error)
	//
	//
	// Inserts row in the database
	InsertMetadata(row util.Row) error
	//
	//
	// Queries the metadata database for selected row. Throws an error if entry doesn't exist
	RetrieveMetadata(row util.Row) (util.Row, error)
	//
	//
	// Select * from table, paged
	ListAllPaged(limit uint, offset uint) ([]util.Row, error)
	//
	//
	// Wrapper for db.Close()
	Close() error
}
