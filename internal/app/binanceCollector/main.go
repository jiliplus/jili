package binancecollector

import (
	"fmt"
	"log"
	"time"

	"github.com/adshao/go-binance"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pelletier/go-toml"
)

const (
	configFile = "binance.toml"
	dbName     = "binance.sqlite3"
)

var (
	client *binance.Client
	db     *gorm.DB
)

func init() {
	// initial client
	config, err := toml.LoadFile(configFile)
	if err != nil {
		msg := fmt.Sprintf("无法导入 %s，%s", configFile, err)
		panic(msg)
	}
	a, s := config.Get("APIKey").(string), config.Get("SecretKey").(string)
	fmt.Printf("APIKey   : %s\n", a)
	fmt.Printf("SecretKey: %s\n", s)
	client = binance.NewClient(a, s)
	fmt.Println("client 初始化完毕")

	// initial db
	db, err = gorm.Open("sqlite3", dbName)
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Printf("%s 数据库已经打开\n", dbName)
}

// Run a binance client to collect historical trades
// NOTICE: 国内的 IP 无法访问 binance 的 API
func Run() {
	defer db.Close()

	rs := newRecords()

	for !rs.isUpdated() {
		symbol, id := rs.first()
		data := request(symbol, id+1)
		save(data)

		last := len(data) - 1
		utc, id := data[last].UTC, data[last].ID

		log.Printf("%s %s", symbol, time.Unix(0, utc*1000000))

		rs.update(utc, id)

	}

	// // 获取历史交易记录
	// res, err := client.NewHistoricalTradesService().Symbol("ETHBTC").FromID(0).Limit(1000).Do(context.TODO())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// r := res[0]
	// fmt.Printf("%d,%d,%s\n", r.ID, r.Time, time.Unix(0, r.Time*1000000))

	// // "bi*" 表示获取所有 bi开头的文件名放入 files
	// files, _ := filepath.Glob("bi*")
	// fmt.Println(files)

	// tp := newTrade("ETHRUB")
	// db.Last(tp)
	// fmt.Println(*tp)

}
