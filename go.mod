module github.com/jujili/jili

go 1.14

require (
	github.com/ThreeDotsLabs/watermill v1.1.1
	github.com/adshao/go-binance v0.0.0-20200414012312-338a1df204bf
	github.com/bearyinnovative/bearychat-go v0.0.0-20181102104846-62b68108f845
	github.com/jinzhu/gorm v1.9.12
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/pelletier/go-toml v1.8.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.5.1
)

replace github.com/jujili/jili/internal => ../internal
