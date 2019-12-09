package beary

import (
	"log"
	"net/http"

	"github.com/pelletier/go-toml"
)

const (
	// ConfigFile 指定了 beary 配置文件的位置
	ConfigFile = "beary.toml"
)

// Channel 是个往 bearychat 发送通知的通道
type Channel struct {
	WebHook string
	From    string
}

// NewChannel make a channel to bearychat.com
func NewChannel() *Channel {
	tree, err := toml.LoadFile(ConfigFile)
	if err != nil {
		log.Fatal("newChannel err:", err)
	}

	c := &Channel{}
	err = tree.Unmarshal(c)
	if err != nil {
		log.Fatal("newChannel tree unmarshal err:", err)
	}

	return c
}

func (c *Channel) send(text, color string) {
	msg := newMessage(c.From, text, color)
	http.Post(c.WebHook, "application/json", msg)
}

func (c *Channel) Verbose(text string) {
	c.send(text, verbose)
}

func (c *Channel) Debug(text string) {
	c.send(text, debug)
}

func (c *Channel) Info(text string) {
	c.send(text, info)
}

func (c *Channel) Warning(text string) {
	c.send(text, warning)
}

func (c *Channel) Fatal(text string) {
	c.send(text, fatal)
}
