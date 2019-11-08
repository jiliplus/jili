package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	filename             = "./test.db"
	wrongCreateStatement = `delete from foo;`
	createStatement      = `
	create table foo (id integer not null primary key, name text);
	delete from foo;
	`
)

func Test_Connect(t *testing.T) {
	ast := assert.New(t)

	// Connect 不存在的filename，使用wrongCreateStatement
	os.Remove(filename)
	_, err := Connect(filename, wrongCreateStatement)
	ast.NotNil(err, "使用错误的创建语句，也还是创建了数据库文件")

	// Connect 不存在的filename，使用createStatement
	os.Remove(filename)
	_, err = Connect(filename, createStatement)
	ast.Nil(err, "文件名和创建表格的语句都对，还出错。")

	// Connect 存在的filename，使用wrongCreateStatement
	defer os.Remove(filename)
	_, err = Connect(filename, wrongCreateStatement)
	ast.Nil(err,
		"打开已经存在的%s，但由于%s，还是报错", filename, wrongCreateStatement)

}
func Test_open(t *testing.T) {
	ast := assert.New(t)

	// open 不存在的filename，使用wrongCreateStatement
	os.Remove(filename)
	_, err := open(filename, wrongCreateStatement)
	ast.NotNil(err, "使用错误的创建语句，也还是创建了数据库文件")

	// open 不存在的filename，使用createStatement
	os.Remove(filename)
	_, err = open(filename, createStatement)
	ast.Nil(err, "文件名和创建表格的语句都对，还出错。")

	// open 存在的filename，使用wrongCreateStatement
	defer os.Remove(filename)
	_, err = open(filename, wrongCreateStatement)
	ast.Nil(err,
		"打开已经存在的%s，但由于%s，还是报错", filename, wrongCreateStatement)

}
func Test_creatDB(t *testing.T) {
	ast := assert.New(t)

	err := createDB(filename, wrongCreateStatement)
	ast.NotNil(err,
		"根据错误的创建语句%s, 也创建了db文件%s", wrongCreateStatement, filename)

	err = createDB(filename, createStatement)
	defer os.Remove(filename)
	ast.Nil(err,
		"在%s文件中，使用%s语句出错", filename, createStatement)
}
