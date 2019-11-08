package db

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetValues(t *testing.T) {
	filename := "testGetValues.db"
	l := 10
	db := prepareDB(filename, l)
	defer os.Remove(filename)

	i := 3
	expectedName := strconv.Itoa(i)
	getValuesStatement := fmt.Sprintf("select name from foo where id = %d", i)

	actualName := ""

	err := db.GetValues(getValuesStatement, &actualName)

	assert.Nil(t, err, "从db获取value数据出错:%s", err)
	assert.Equal(t, expectedName, actualName, "从db中获取的name不对")
}
