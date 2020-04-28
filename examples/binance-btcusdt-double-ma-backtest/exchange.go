package main

import (
	"context"
	"log"
	"math"

	"github.com/ThreeDotsLabs/watermill/message"
)

type exchange struct {
	coin, money *asset
	done        context.CancelFunc
}

type asset struct {
	free, frozen, total float64
}

func newAsset(initailFree float64) *asset {
	return &asset{
		free: initailFree,
	}
}

type tick struct {
}

type pubsub interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(topic string) (<-chan *message.Message, error)
	Close() error
}

type order struct {
	ID int64
	// TODO: 把 OrderType 的 string 转变成枚举
	OrderType       string
	Price, Quantity float64
	Next, Prev      *order
}

func newOrder(id int64, orderType string, price, quantity float64) *order {
	return &order{
		ID:        id,
		OrderType: orderType,
		Price:     price,
		Quantity:  quantity,
	}
}

type orders struct {
	head, tail *order
}

func newOrders() *orders {
	head := newOrder(math.MinInt64, "", -math.MaxFloat64, 0)
	tail := newOrder(math.MaxInt64, "", math.MaxFloat64, 0)
	head.Next = tail
	tail.Prev = head
	return &orders{head: head, tail: tail}
}

func (o *orders) add(order *order) {
	curr, next := o.head, o.head.Next
	for order.Price > next.Price {
		curr, next = next, next.Next
	}
	// now,
	// order.Price <= next.Price
	curr.Next, order.Next = order, next
	next.Prev, order.Prev = order, curr
}

func (o *orders) cancel(id int64) {
	curr := o.head.Next
	for curr.ID != id {
		curr = curr.Next
	}
	if curr == o.tail {
		return
	}
	curr.Next.Prev, curr.Prev.Next = curr.Prev, curr.Next
}

func (o *orders) cancelAll() {
	o.head.Next = o.tail
	o.tail.Prev = o.head
}

func (o *orders) () {
	o.head.Next = o.tail
	o.tail.Prev = o.head
}



func newExchange(ctx context.Context, source <-chan tick, pubsub pubsub, initialMoney float64) *exchange {
	exCtx, cancel := context.WithCancel(ctx)
	ex := &exchange{
		coin:  newAsset(0),
		money: newAsset(initialMoney),
		done:  cancel,
	}

	// 跑一个虚拟的盘口

	go func() {
		getOrder, err := pubsub.Subscribe("order")
		if err != nil {
			panic("pubsub.Subscribe(order) error" + err.Error())
		}
		orders := newOrders()
		for {
			select {
			case <-exCtx.Done():
				log.Println("ctx done.")
			case <-getOrder:
			}
		}
	}()

	return ex
}
