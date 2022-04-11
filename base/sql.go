package base

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type SqlDB struct {
	db     *sql.DB
	logger *logrus.Logger
	table  string
}

func NewSqlDatabase(logger *logrus.Logger, table string) DB {
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
	err = db.Ping()
	if err != nil {
		return err
	}

	sqldb.db = db

	return nil
}

func (sqldb *SqlDB) Init() error {
	sqldb.logger.Debug("Creating table, if doesn't exist")
	statementString := "CREATE TABLE IF NOT EXISTS " + sqldb.table + " ( uuid uuid PRIMARY KEY, fileName varchar(255) NOT NULL);"
	sqldb.logger.Debug(statementString)
	statement, err := sqldb.db.Prepare(statementString)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	sqldb.logger.Debug("Create Table query executed")

	// sqldb.logger.Debug("Inserting test row in table")
	// statementString = "INSERT INTO " + sqldb.table + " VALUES( '" + uuid.New().String() + "', 'name' );"
	// sqldb.logger.Debug(statementString)
	// if statement, err = sqldb.db.Prepare(statementString); err != nil {
	// 	return err
	// }
	// if _, err = statement.Exec(); err != nil {
	// 	return err
	// }
	// sqldb.logger.Debug("Insert query executed")

	return nil
}

func (sqldb *SqlDB) InsertMetadata(row Row) error {
	statementString := "INSERT INTO " + sqldb.table + " VALUES( $1, $2 );"
	sqldb.logger.Debug(statementString)
	statement, err := sqldb.db.Prepare(statementString)
	if err != nil {
		return err
	}
	sqldb.logger.Debug("Insert query prepared")
	res, err := statement.Exec(row.Uuid, row.FileName)
	if err != nil {
		return err
	}
	sqldb.logger.Debug("Insert query executed")
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	sqldb.logger.Debugf("Created %s rows", rowCnt)

	return nil
}

func (sqldb *SqlDB) Close() error {
	return sqldb.db.Close()
}
