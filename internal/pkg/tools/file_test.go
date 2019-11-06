package tools

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsExist(t *testing.T) {
	Convey("假如生成了临时文件", t, func() {
		tmp, err := ioutil.TempFile("", "TestIsExist")
		So(err, ShouldBeNil)
		Convey("当检测其是否存在时，", func() {
			Convey("应该返回 true", func() {
				So(IsExist(tmp.Name()), ShouldBeTrue)
			})
		})
		os.Remove(tmp.Name())
		Convey("当删除文件后,", func() {
			Convey("应该返回 false", func() {
				So(IsExist(tmp.Name()), ShouldBeFalse)
			})
		})
	})
}
