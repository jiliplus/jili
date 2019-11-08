package db

import (
	"database/sql"
	"sync"

	"github.com/aQuaYi/GoKit"
)

// Connect 返回一个数据库对象
func Connect(filename string, createStatement string) (DBer, error) {
	db, err := open(filename, createStatement)
	if err != nil {
		return nil, err
	}

	return &DB{
		name: filename,
		DB:   db,
	}, nil
}

var mutex sync.Mutex

// open 链接上了数据库
func open(filename string, createStatement string) (*sql.DB, error) {
	{ //为了不重复创建数据库，加个锁
		mutex.Lock()
		defer mutex.Unlock()

		//如果不存在数据库文件不存在，就创建一个新的
		if !GoKit.Exist(filename) {
			if err := createDB(filename, createStatement); err != nil {
				return nil, err
			}
		}
	}

	// 查看sql.Open的源码可知，只有在驱动错误的情况下，才会报错。
	// NOTICE: 如果以后，可以自由选择数据库驱动的时候，要检查err
	db, _ := sql.Open("sqlite3", filename)

	return db, nil
}

func createDB(filename string, createStatement string) error {
	// 查看sql.Open的源码可知，只有在驱动错误的情况下，才会报错。
	// NOTICE: 如果以后，可以自由选择数据库驱动的时候，要检查err
	db, _ := sql.Open("sqlite3", filename)
	defer db.Close()

	_, err := db.Exec(createStatement) //数据库执行创建语句
	if err != nil {
		return GoKit.Err(err,
			`在%s数据库中执行语句"%s"失败`, filename, createStatement)
	}

	return nil
}
