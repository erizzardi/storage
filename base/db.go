package base

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
	// Inserts row in the database
	//----------------------------
	InsertMetadata(row Row) error
	//
	// Wrapper for db.Close()
	//-----------------------
	Close() error
}
