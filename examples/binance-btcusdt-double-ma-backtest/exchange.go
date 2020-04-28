package main

import "github.com/ThreeDotsLabs/watermill/message"

type exchange struct {
	coin, money asset
}

type asset struct {
	free, frozen, total float64
}

type tick struct {
}

type pubsub interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(topic string) (<-chan *message.Message, error)
	Close() error
}

func newExchange(source <-chan tick, pubsub pubsub, initialMoney float64) {

	return
}
