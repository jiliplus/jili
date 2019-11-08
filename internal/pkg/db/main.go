// Package db to manage database
package db

import (
	"database/sql"
)

//DBer 是定制数据库的接口
type DBer interface {
	Name() string
	Insert(statement string, data interface{}) error
	GetRows(statement string, structPtr interface{}) ([]interface{}, error)
	GetValues(statement string, dest ...interface{}) error
}

//DB 定制的sql数据库
type DB struct {
	name string
	*sql.DB
}

//Name 返回数据库的名称，也是数据库的存放地址
func (db *DB) Name() string {
	return db.name
}
