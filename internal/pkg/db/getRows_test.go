package db

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareDB(filename string, length int) DBer {
	db, _ := Connect(filename, createStatement)

	tds := makeTDS(length)

	db.Insert(insertStatement, tds)

	return db
}

func Test_GetRows(t *testing.T) {
	filename := "testGetRows.db"
	l := 10
	db := prepareDB(filename, l)
	defer os.Remove(filename)

	getRowsStatement := "select id, name from foo"
	data, err := db.GetRows(getRowsStatement, &testData{})
	assert.Nil(t, err, "从db获取Rows数据出错:%s", err)

	for i := 0; i < l; i++ {
		d, ok := data[i].(testData)
		assert.True(t, ok, "从数据库获取的数据，类型不对")
		assert.Equal(t, int64(i), d.ID, "从数据库获取的数据，顺序不对")
		assert.Equal(t, strconv.Itoa(i), d.Name, "从数据库获取的数据，Name不对")
	}
}

func Test_makeARow_structPointer(t *testing.T) {
	ast := assert.New(t)
	td := testData{
		ID:   1,
		Name: "2",
	}
	is, s := makeRow(&td)
	ast.Equal(td, s.Interface(), "转换后的td不同了")
	ast.Equal(&td.ID, is[0], "转换后的td.ID不同了")
	ast.Equal(&td.Name, is[1], "转换后的td.Name不同了")
}
