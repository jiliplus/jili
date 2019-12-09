package main

import (
	"net/http"

	"github.com/bearyinnovative/bearychat-go"
)

func main() {
	m := bearychat.Incoming{
		Text:         "Hello, **BearyChat**",
		Markdown:     true,
		Notification: "Hello, BearyChat in Notification",
	}

	output, _ := m.Build()

	http.Post("https://hook.bearychat.com/=bwGvD/incoming/98272d49d0be2614b69670977be7403e", "application/json", output)
}
