package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Product is ...
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	// file::memory:?cache=shared
	db1, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic("failed to connect database")
	}
	defer db1.Close()
	db2, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic("failed to connect database")
	}
	defer db2.Close()

	// Migrate the schema
	db1.AutoMigrate(&Product{})
	db2.AutoMigrate(&Product{})

	// 创建
	db1.Create(&Product{Code: "L1212", Price: 1000})

	// 读取
	var product Product
	fmt.Println("original: ", product)
	db2.First(&product, 1) // 查询id为1的product
	// 可以看到从 db2 中查询到了 插入到 db1 中的数据
	fmt.Println("from db2: ", product)
}
