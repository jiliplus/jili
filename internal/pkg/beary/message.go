package beary

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bearyinnovative/bearychat-go"
)

// 5 个等级对应的颜色
const (
	verbose = "#968387"
	debug   = "#0C50A3"
	info    = "#80C31C"
	warning = "#FFBF01"
	fatal   = "#FF4401"
)

func now() string {
	return time.Now().Format("2006-01-02 15:04:05.999")
}

func newMessage(from, text, color string) io.Reader {
	incoming := bearychat.Incoming{
		Text:        fmt.Sprintf("**%s** %s ", from, now()),
		Attachments: []bearychat.IncomingAttachment{newAttachment(text, color)},
	}
	res, err := incoming.Build()
	if err != nil {

		log.Fatalf("以下内容，无法发送到 BearyChat，From:%s, Text:%s", from, text)
	}
	return res
}

func newAttachment(text, color string) bearychat.IncomingAttachment {
	return bearychat.IncomingAttachment{
		Text:  text,
		Color: color,
	}
}
