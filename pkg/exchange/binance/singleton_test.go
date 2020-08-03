package binance

import (
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/jujili/jili/pkg/tools"
)

func TestExit(t *testing.T) {
	Convey("假设 ConfigFile 不存在", t, func() {
		Convey("在退出前", func() {
			Convey("ConfigFile 不存在", func() {
				So(tools.IsExist(ConfigFile), ShouldBeFalse)
			})
		})
		Convey("在退出时", func() {
			Convey("会 panic", func() {
				So(func() {
					exit()
				}, ShouldPanicWith,
					fmt.Sprintf("在当前目录没有找到 %s 文件，已帮你生成一个空的 %s，请填写完成后，再启动程序。", ConfigFile, ConfigFile))
			})
		})
		Convey("在退出后", func() {
			Convey("configFile 存在", func() {
				So(tools.IsExist(ConfigFile), ShouldBeTrue)
			})
			Reset(func() {
				os.Remove(ConfigFile)
			})
		})
	})
}
