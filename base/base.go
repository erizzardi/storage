package base

import "database/sql"

type Database interface {
	List(table string)
	Insert(table string, key string)
}

type SqlDatabase struct{ db *sql.DB }

func NewSqlDatabase(db *sql.DB) *SqlDatabase { return &SqlDatabase{db} }

func (db *SqlDatabase) List(table string) {
	/**   **/
}

func (db *SqlDatabase) Insert(table string, key string) {
	/**   **/
}
