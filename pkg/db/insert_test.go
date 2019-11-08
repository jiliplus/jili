package db

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	ID   int64
	Name string
}

var (
	insertStatement      = "insert into foo(id, name) values(?,?)"
	wrongInsertStatement = "insert into foo(id, name, age) values(?,?,?)"
)

//Attributes 实现了database.Attributer接口

func makeTDS(length int) []testData {
	res := make([]testData, length)

	for i := 0; i < length; i++ {
		res[i] = testData{
			ID:   int64(i),
			Name: strconv.Itoa(i),
		}
	}

	return res
}

func Test_Insert_wrongInsertStatement(t *testing.T) {
	filename := "./insert_test.db"
	db, _ := Connect(filename, createStatement)
	defer os.Remove(filename)
	l := 20
	tds := makeTDS(l)

	assert.NotNil(t, db.Insert(wrongInsertStatement, tds), "使用了错误的插入语句，但是没有报错")
}
func Test_Insert_AfterCloseDB(t *testing.T) {
	filename := "./insert_test.db"
	db, _ := Connect(filename, createStatement)
	defer os.Remove(filename)
	l := 20
	tds := makeTDS(l)

	go func() {
		base, _ := db.(*DB)
		base.DB.Close()
	}()
	time.Sleep(100 * time.Millisecond)
	err := db.Insert(insertStatement, tds)
	assert.NotNil(t, err, "插入已被关闭的数据库，也没有报错。")
}
func Test_Insert_badData(t *testing.T) {
	filename := "./insert_test.db"
	db, _ := Connect(filename, createStatement)
	defer os.Remove(filename)
	l := 20

	badTDS := make([]testData, l)
	for i := 0; i < l; i++ {
		badTDS[i] = testData{
			// NOTICE:  我设置了重复的主键
			ID:   0,
			Name: strconv.Itoa(i),
		}
	}

	err := db.Insert(insertStatement, badTDS)
	assert.NotNil(t, err, "插入错误的数据时，没有报错。")
}

func Test_Insert(t *testing.T) {
	filename := "./insert_test.db"
	db, _ := Connect(filename, createStatement)
	defer os.Remove(filename)
	l := 20
	tds := makeTDS(l)

	err := db.Insert(insertStatement, tds)
	assert.Nil(t, err, "插入正确的数据时，也会报错。")
}

func Test_makeIS_NoSlice(t *testing.T) {
	i := 1
	is := makeIS(i)

	assert.Equal(t, i, is[0], "单个输入，没能变成切片")
}

func Test_makeIS_Slice(t *testing.T) {
	is := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	iss := makeIS(is)

	for i, v := range is {
		assert.Equal(t, v, iss[i], "转换后，第%d个元素变了", i)
	}
}
