package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xormplus/xorm"
)

var engine *xorm.Engine

func main() {
	var err error
	engine, err = xorm.NewSqlite3(":memory:")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Ping error: ", engine.Ping())
	fmt.Print("has \"ha\" table")
	fmt.Println(engine.IsTableExist("ha"))

	fmt.Print("DBMates:")
	fmt.Println(engine.DBMetas())

	u := user{ID: 1}

	err = engine.CreateTables(u)
	if err != nil {
		fmt.Println("create tables err:", err)
	}

	tables, _ := engine.DBMetas()
	fmt.Print("tables:")
	fmt.Println(tables[0].Name)

	engine.Table("user").Insert(&user{ID: 2})
	engine.Table("user").Insert(&user{ID: 3})
	engine.Table("user").Insert(&user{ID: 4})

	isEmpty, err := engine.IsTableEmpty(u)
	if err != nil {
		fmt.Println("is table empty", err)
	} else {
		fmt.Println("is table user empty:", isEmpty)
	}

	u2 := new(user)
	rows, err := engine.Table("user").Rows(u2)
	if err != nil {
		fmt.Println("engine row err:", err)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(u2)
		if err != nil {
			fmt.Println("rows scan err:", err)
		}
		fmt.Println("u2", u2)
	}

	u.ID = 0
	b, err := engine.Table("user").Desc("id").Get(&u)
	fmt.Println("b,err:", b, err)

	fmt.Println("u is", u)
}

type user struct {
	ID int64 `xorm:"pk 'id'"`
}

func (*user) TableName() string {
	return "user"
}
