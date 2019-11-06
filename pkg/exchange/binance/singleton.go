package binance

import "sync"

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
	once.Do(func() {
		instance = &Client{}
	})
	return instance
}
