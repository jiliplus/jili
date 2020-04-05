package main

import (
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
)

func main() {
	msg := message.NewMessage("1", []byte("foo"))

	go func() {
		m := msg.Copy()
		time.Sleep(time.Millisecond * 10)
		m.Ack()
	}()

	select {
	case <-msg.Acked():
		log.Print("ack received")
	case <-msg.Nacked():
		log.Print("nack received")
	case <-time.After(time.Second):
		log.Print("timeout")
	}

}
