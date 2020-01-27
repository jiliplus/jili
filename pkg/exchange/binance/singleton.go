package binance

import (
	"fmt"
	"sync"

	"github.com/pelletier/go-toml"

	"github.com/aQuaYi/jili/internal/pkg/tools"
)

const (
	// ConfigFile contains all config of binance
	ConfigFile = "./binance.toml"
)

// Client is ...
type Client struct {
	APIKey    string
	SecretKey string

	BaseURL string
}

var instance *Client
var once sync.Once

// New will return a new client
// TODO: 修改此处返回值的类型
func New() *Client {
	if !tools.IsExist(ConfigFile) {
		exit()
	}
	once.Do(func() {
		tree, _ := toml.LoadFile(ConfigFile)
		err := tree.Unmarshal(instance)
		if err != nil {
			msg := fmt.Sprintf("读取 %s 文件错误，请核查。", ConfigFile)
			panic(msg)
		}
	})
	return instance
}

func exit() {
	c := &Client{}
	data, _ := toml.Marshal(c)
	tools.Create(ConfigFile, data)
	msg := fmt.Sprintf("在当前目录没有找到 %s 文件，已帮你生成一个空的 %s，请填写完成后，再启动程序。", ConfigFile, ConfigFile)
	panic(msg)
}
