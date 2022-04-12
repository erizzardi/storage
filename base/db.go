package base

import "database/sql"

// Interface for database connection and operations.
// Implementation for sql database: sql.go
type DB interface {
	//
	// Enstablishes connection to the database
	//----------------------------------------
	Connect(driver string, dsl string) error
	//
	// Creates the table in the database
	//----------------------------------
	Init() error
	//
	// Prepares and executes 'statement' operation - relies on db.Exec()
	//------------------------------------------------------------------
	Exec(statement string) (sql.Result, error)
	//
	// Inserts row in the database
	//----------------------------
	InsertMetadata(row Row) error
	//
	// Queries the database for selected row. Throws an error if entry doesn't exist
	//------------------------------------------------------------------------------
	RetrieveMetadata(row Row) (Row, error)
	//
	// Wrapper for db.Close()
	//-----------------------
	Close() error
}
