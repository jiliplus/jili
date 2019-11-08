package db

import (
	"errors"
	"fmt"
)

// GetValues 获取查询语句的多个值
// 在dest存放变量的指针
// NOTICE: dest中变量的指针的顺序，需要与statement中查询的一样
func (db *DB) GetValues(statement string, dest ...interface{}) error {
	stmt, err := db.Prepare(statement)
	if err != nil {
		msg := fmt.Sprintf("对%s使用以下语句查询\n%s\n出现错误:%s", db.Name(), statement, err)
		return errors.New(msg)
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(dest...) //NOTICE: Scan的参数必须打上...
	if err != nil {
		msg := fmt.Sprintf("database.GetValues：对%s查询%s出来的值，Scan完毕后，出错:%s", db.Name(), statement, err)
		return errors.New(msg)
	}

	return nil
}
