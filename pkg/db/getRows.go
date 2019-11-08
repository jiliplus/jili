package db

import (
	"errors"
	"fmt"
	"reflect"
)

//GetRows 描述了从数据库中查询的过程
func (db *DB) GetRows(queryStatement string, StructPtr interface{}) ([]interface{}, error) {
	rows, err := db.Query(queryStatement)
	if err != nil {
		msg := fmt.Sprintf("对%s使用以下语句查询\n%s\n出现错误:%s", db.Name(), queryStatement, err)
		return nil, errors.New(msg)
	}
	defer rows.Close()

	result := []interface{}{}
	for rows.Next() {
		row, s := makeRow(StructPtr)
		err := rows.Scan(row...)
		if err != nil {
			msg := fmt.Sprintf("对%s查询%s出来的rows进行Scan时，出错:%s", db.Name(), queryStatement, err)
			return nil, errors.New(msg)
		}
		result = append(result, s.Interface())
	}

	err = rows.Err()
	if err != nil {
		msg := fmt.Sprintf("对%s查询%s出来的rows，Scan完毕后，出错:%s", db.Name(), queryStatement, err)
		return nil, errors.New(msg)
	}

	return result, nil
}

func makeRow(structPtr interface{}) ([]interface{}, reflect.Value) {
	v := reflect.ValueOf(structPtr).Elem()
	leng := v.NumField()
	onerow := make([]interface{}, leng)
	for i := 0; i < leng; i++ {
		onerow[i] = v.Field(i).Addr().Interface()
	}

	return onerow, v
}
