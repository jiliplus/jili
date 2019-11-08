package db

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aQuaYi/GoKit"
)

//Insert 描述了向DB内插入数据的过程
func (db *DB) Insert(insertStatement string, data interface{}) error {
	//启动insert事务
	transaction, err := db.Begin()
	if err != nil {
		return GoKit.Err(err, "%s无法启动一个insert事务", db.Name())
	}
	defer transaction.Commit()

	//为insert事务进行准备工作
	stmt, err := transaction.Prepare(insertStatement)
	if err != nil {
		msg := fmt.Sprintf("%s的insert事务的准备以下insert语句时失败\n%s\n失败原因: %s", db.Name(), insertStatement, err)
		return errors.New(msg)
	}
	defer stmt.Close()

	dataSlice := makeIS(data)
	//按行插入
	for _, d := range dataSlice {
		_, err := stmt.Exec(attributes(d)...)
		if err != nil {
			attrs := fmt.Sprint(d)
			msg := fmt.Sprintf("%s在插入%s出错: %s", db.Name(), attrs, err)
			//NOTICE: 经过再三的思考，我决定在插入出错后，不要直接关闭程序。由程序的调用方来决定，如何处理错误。
			return errors.New(msg)
		}
	}

	return nil
}

func attributes(struc interface{}) []interface{} {
	t := reflect.TypeOf(struc)
	v := reflect.ValueOf(struc)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	n := v.NumField()
	res := make([]interface{}, n)
	for i := 0; i < n; i++ {
		res[i] = v.Field(i).Interface()
	}

	return res
}

func makeIS(itf interface{}) []interface{} {
	t := reflect.TypeOf(itf)
	if t.Kind() != reflect.Slice {
		res := make([]interface{}, 1)
		res[0] = itf
		return res
	}

	v := reflect.ValueOf(itf)
	l := v.Len()

	res := make([]interface{}, l)
	for i := 0; i < l; i++ {
		res[i] = v.Index(i).Interface()
	}
	return res
}
