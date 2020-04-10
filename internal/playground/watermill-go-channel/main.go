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
	pubsub := gochannel.NewGoChannel(
		gochannel.Config{
			// OutputChannelBuffer:            10,
			BlockPublishUntilSubscriberAck: true,
		},
		watermill.NewStdLogger(false, false),
	)

	// messages, err := pubsub.Subscribe(context.Background(), TOPIC)
	// if err != nil {
	// 	panic(err)
	// }
	// go process(messages)

	go subscribe(1, pubsub)
	go subscribe2(2, pubsub)

	publishMessages(pubsub)

	// pubsub.Close()
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
}

// func process(messages <-chan *message.Message) {
// 	for msg := range messages {
// 		log.Printf("received message\t: %s, payload: %s\n", msg.UUID, string(msg.Payload))
// 		// we need to Acknowledge that we received and processed the message,
// 		// otherwise, it will be resent over and over again.
// 		msg.Ack()
// 	}
// }

func subscribe(id int, pub message.Subscriber) {
	messages, err := pub.Subscribe(context.Background(), TOPIC)
	if err != nil {
		panic(err)
	}
	for msg := range messages {
		log.Printf("%d,received message\t: %s, payload: %s\n", id, msg.UUID, string(msg.Payload))
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}

func subscribe2(id int, pub message.Subscriber) {
	messages, err := pub.Subscribe(context.Background(), TOPIC)
	if err != nil {
		panic(err)
	}
	for msg := range messages {
		log.Printf("%d,received message\t: %s, payload: %s\n", id, msg.UUID, string(msg.Payload))
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
