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

func NewSqlDatabase(db *sql.DB, logger *util.Logger, table string) DB {
	return &SqlDB{
		db:     db,
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

func (sqldb *SqlDB) InsertMetadata(row util.Row) error {
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
	sqldb.logger.Debugf("Created %d rows", rowCnt)

	return nil
}

func (sqldb *SqlDB) RetrieveMetadata(row util.Row) (util.Row, error) {
	var ret util.Row

	statementString := "SELECT * FROM " + sqldb.table + " WHERE uuid = $1;"
	sqldb.logger.Debug(statementString)
	rows, err := sqldb.Query(statementString, row.Uuid)
	if err != nil {
		return util.Row{}, err
	}
	// statement, err := sqldb.db.Prepare(statementString)
	// if err != nil {
	// 	return util.Row{}, err
	// }
	// sqldb.logger.Debug("Select query prepared")
	// rows, err := statement.Query(row.Uuid)
	// if err != nil {
	// 	return util.Row{}, err
	// }
	sqldb.logger.Debug("Insert query executed")
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

	return ret, nil
}

func (sqldb *SqlDB) ListAllPaged(limit uint, offset uint) ([]util.Row, error) {
	ret := make([]util.Row, 0)
	var tempUuid, tempFileName string

	statementString := "SELECT * FROM " + sqldb.table + " LIMIT $1 OFFSET $2;"
	sqldb.logger.Debug(statementString)
	rows, err := sqldb.Query(statementString, limit, offset)
	if err != nil {
		return []util.Row{}, err
	}
	sqldb.logger.Debug("Insert query executed")
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
