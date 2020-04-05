// Sources for https://watermill.io/docs/getting-started/
package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

// TOPIC 定义了双方的话题
const TOPIC = "example.topic"

func main() {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	messages, err := pubSub.Subscribe(context.Background(), TOPIC)
	if err != nil {
		panic(err)
	}

	go publishMessages(pubSub)

	process(messages)
}

func publishMessages(publisher message.Publisher) {
	for i := 0; i < 10; i++ {
		msg := message.NewMessage(strconv.Itoa(i), []byte("Hi, world!"))

		if err := publisher.Publish(TOPIC, msg); err != nil {
			panic(err)
		}
		log.Printf("\tsended message\t: %s, payload: %s\n", msg.UUID, string(msg.Payload))

		time.Sleep(time.Second * 1)

	}
	publisher.Close()
}

func process(messages <-chan *message.Message) {
	for msg := range messages {

		log.Printf("received message\t: %s, payload: %s\n", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
