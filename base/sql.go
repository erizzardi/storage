package base

import (
	"database/sql"

	"github.com/erizzardi/storage/util"
	_ "github.com/lib/pq"
)

type SqlDB struct {
	db     *sql.DB
	logger *util.Logger
	table  string
}

func NewSqlDatabase(logger *util.Logger, table string) DB {
	return &SqlDB{
		db:     &sql.DB{},
		logger: logger,
		table:  table,
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

func (sqldb *SqlDB) Init() error {
	sqldb.logger.Debug("Creating 'meta' table, if doesn't exist")
	statementString := "CREATE TABLE IF NOT EXISTS " + sqldb.table + " ( uuid uuid PRIMARY KEY, fileName varchar(255) NOT NULL );"
	if _, err := sqldb.Exec(statementString); err != nil {
		return err
	}
	sqldb.logger.Debug("Create Table 'meta' query executed")

	sqldb.logger.Debug("Creating 'bucket' table, if doesn't exist")
	statementString = "CREATE TABLE IF NOT EXISTS " + sqldb.table + " ( name varchar(255) PRIMARY KEY, owner varchar(255) NOT NULL);"
	if _, err := sqldb.Exec(statementString); err != nil {
		return err
	}

	sqldb.logger.Debug("Create Table 'bucket' query executed")

	return nil
}

func (sqldb *SqlDB) InsertMetadata(row Row) error {
	statementString := "INSERT INTO " + sqldb.table + " VALUES( $1, $2 );"
	sqldb.logger.Debug(statementString)

	res, err := sqldb.Exec(statementString, row.Uuid, row.FileName)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	sqldb.logger.Debugf("Created %s rows", rowCnt)

	return nil
}

func (sqldb *SqlDB) RetrieveMetadata(row Row) (Row, error) {
	var ret Row

	statementString := "SELECT * FROM " + sqldb.table + " WHERE uuid = $1;"
	sqldb.logger.Debug(statementString)
	statement, err := sqldb.db.Prepare(statementString)
	if err != nil {
		return Row{}, err
	}
	sqldb.logger.Debug("Select query prepared")
	rows, err := statement.Query(row.Uuid)
	if err != nil {
		return Row{}, err
	}
	sqldb.logger.Debug("Insert query executed")
	for rows.Next() {
		err := rows.Scan(&ret.Uuid, &ret.FileName)
		if err != nil {
			return Row{}, err
		}
		sqldb.logger.Debugf("Row read. Retrieved %+v\n", ret)
	}
	err = rows.Err()
	if err != nil {
		return Row{}, err
	}
	sqldb.logger.Debug("Row scanning ended")

	return ret, nil

}

func (sqldb *SqlDB) Close() error {
	return sqldb.db.Close()
}
