package base

import (
	"database/sql"
	"fmt"

	"github.com/erizzardi/storage/util"
	_ "github.com/lib/pq"
)

type SqlDB struct {
	db     *sql.DB
	logger *util.Logger
	tables []table
}

func NewSqlDatabase(logger *util.Logger) DB {
	return &SqlDB{
		db:     &sql.DB{},
		logger: logger,
		//-------------------
		// Database Structure
		//-------------------
		tables: []table{
			{
				name: "meta",
				columns: []column{
					newColumn("uuid", "uuid", true, false),
					newColumn("fileName", "varchar(255)", false, true),
				},
				labels: map[string]any{
					"content": "metadata",
				},
			},
			{
				name: "bucket",
				columns: []column{
					newColumn("name", "varchar(255)", true, false),
					newColumn("owner", "varchar(255)", false, true),
				},
				labels: map[string]any{
					"content": "bucket",
				},
			},
		},
	}
}

func (sqldb *SqlDB) Connect(driver string, dsn string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	sqldb.logger.Debug("Pinging database...")
	err = db.Ping()
	if err != nil {
		return err
	}

	sqldb.db = db

	return nil
}

func (sqldb *SqlDB) Exec(statementString string, params ...any) (sql.Result, error) {
	sqldb.logger.Debug(statementString)
	statement, err := sqldb.db.Prepare(statementString)
	if err != nil {
		return nil, err
	}
	res, err := statement.Exec(params...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (sqldb *SqlDB) Query(statementString string, params ...any) (*sql.Rows, error) {
	sqldb.logger.Debug(statementString)
	statement, err := sqldb.db.Prepare(statementString)
	if err != nil {
		return nil, err
	}
	sqldb.logger.Debug("Select query prepared")
	rows, err := statement.Query(params...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Init() loads structure into database
func (sqldb *SqlDB) Init() error {

	for _, table := range sqldb.tables {
		sqldb.logger.Debugf("Creating '%s' table, if doesn't exist", table.name)
		statementString := "CREATE TABLE IF NOT EXISTS " + table.name + " ("
		for i, column := range table.columns {
			if i != 0 {
				statementString += ","
			}
			statementString += column.toString()
		}
		statementString += ")"
		sqldb.logger.Debug("Table creation statement: " + statementString)
		if _, err := sqldb.Exec(statementString); err != nil {
			return err
		}
		sqldb.logger.Debugf("Create Table '%s' executed", table.name)
	}
	return nil
}

// tearDown() drops all the tables created by Init(). To be used in tests! Thus, unexported.
func (sqldb *SqlDB) tearDown() error {

	for _, table := range sqldb.tables {
		sqldb.logger.Debugf("Dropping table '%s'", table.name)
		statementString := "DROP TABLE " + table.name
		sqldb.logger.Debug("Table creation statement: " + statementString)
		if _, err := sqldb.Exec(statementString); err != nil {
			return err
		}
		sqldb.logger.Debugf("Dropped table '%s'", table.name)
	}
	return nil
}

func (sqldb *SqlDB) InsertMetadata(row util.Row) error {

	statementString := "INSERT INTO " + sqldb.GetTableFromLabel("metadata") + " VALUES( $1, $2 );"
	sqldb.logger.Debug(statementString)

	res, err := sqldb.Exec(statementString, row.Uuid, row.FileName)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	sqldb.logger.Debugf("Created %d row/s", rowCnt)

	return nil
}

func (sqldb *SqlDB) RetrieveMetadata(key, value string) (util.Row, error) {

	var ret util.Row

	// POSSIBLE SQL INJECTION
	statementString := "SELECT * FROM " + sqldb.GetTableFromLabel("metadata") + " WHERE " + key + " = $1;"
	// sqldb.logger.Debug(statementString)

	rows, err := sqldb.Query(statementString, value)
	if err != nil {
		return util.Row{}, err
	}
	sqldb.logger.Debug("Select query executed")
	for rows.Next() {
		err := rows.Scan(&ret.Uuid, &ret.FileName)
		if err != nil {
			return util.Row{}, err
		}
		sqldb.logger.Debugf("Row read. Retrieved %+v\n", ret)
	}
	err = rows.Err()
	if err != nil {
		return util.Row{}, err
	}
	sqldb.logger.Debug("Row scanning ended")

	// Checks for empty response
	if (ret == util.Row{}) {
		return ret, NotFoundError
	}

	return ret, nil
}

func (sqldb *SqlDB) DeleteMetadata(key, value string) error {

	// POSSIBLE SQL INJECTION
	statementString := "DELETE FROM " + sqldb.GetTableFromLabel("metadata") + " WHERE " + key + " = $1;"
	sqldb.logger.Debug(statementString)
	res, err := sqldb.Exec(statementString, value)
	if err != nil {
		return err
	}

	if rowCnt, err := res.RowsAffected(); err != nil {
		return util.InternalServerError{Message: err.Error()}
	} else if rowCnt > 1 {
		sqldb.logger.Errorf("DELETE operation affected %d lines.", rowCnt)
		return util.InternalServerError{Message: fmt.Sprintf("DELETE operation affected %d lines.", rowCnt)}
	} else if rowCnt == 0 {
		sqldb.logger.Errorf("Error: file %s non existing.", value)
		return util.BadRequestError{Message: ("file " + value + "does not exist.")}
	}

	return nil
}

func (sqldb *SqlDB) ListAllPaged(limit uint, offset uint) ([]util.Row, error) {

	ret := make([]util.Row, 0)
	var tempUuid, tempFileName string

	statementString := "SELECT * FROM " + sqldb.GetTableFromLabel("metadata") + " LIMIT $1 OFFSET $2;"
	sqldb.logger.Debug(statementString)
	rows, err := sqldb.Query(statementString, limit, offset)
	if err != nil {
		return []util.Row{}, err
	}
	sqldb.logger.Debug("Select query executed")
	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&tempUuid, &tempFileName)
		if err != nil {
			return []util.Row{}, err
		}
		ret = append(ret, util.Row{Uuid: tempUuid, FileName: tempFileName})
		sqldb.logger.Debugf("Row read. Retrieved %+v\n", ret)
	}
	err = rows.Err()
	if err != nil {
		return []util.Row{}, err
	}
	sqldb.logger.Debug("Row scanning ended")

	return ret, nil
}

func (sqldb *SqlDB) Close() error {
	return sqldb.db.Close()
}

//============
// Miscellanea
//============
func (sqldb *SqlDB) GetTableFromLabel(label string) string {

	var tableName string
	for _, table := range sqldb.tables {
		if table.labels["content"] == "metadata" {
			tableName = table.name
		}
	}
	return tableName
}
