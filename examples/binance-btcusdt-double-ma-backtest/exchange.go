package main

import (
	"bytes"
	"context"
	"encoding/gob"
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
	Type            string
	Price, Quantity float64
	Next            *order
}

func newOrder(id int64, orderType string, price, quantity float64) *order {
	return &order{
		ID:       id,
		Type:     orderType,
		Price:    price,
		Quantity: quantity,
	}
}

// TODO: 利用有限队列实现 sellOrders

type sellOrders struct {
	head *order
}

func newSellOrders() *sellOrders {
	head := newOrder(math.MinInt64, "", -math.MaxFloat64, 0)
	return &sellOrders{head: head}
}

func (o *sellOrders) push(order *order) {
	curr, next := o.head, o.head.Next
	for next != nil && order.Price > next.Price {
		curr, next = next, next.Next
	}
	// now,
	// order.Price <= next.Price
	curr.Next, order.Next = order, next
}

func (o *sellOrders) pop() *order {
	if o.head.Next == nil {
		return nil
	}
	res := o.head.Next
	o.head.Next = o.head.Next.Next
	res.Next = nil
	return res
}

func (o *sellOrders) cancel(id int64) {
	curr := o.head
	for curr.Next != nil && curr.Next.ID != id {
		curr = curr.Next
	}
	if curr.Next == nil {
		return
	}
	curr.Next = curr.Next.Next
}

func (o *sellOrders) cancelAll() {
	o.head.Next = nil
}

func (o *sellOrders) canSell(price float64) bool {
	return o.head.Next.Price <= price
}

type buyOrders struct {
	head *order
}

func newBuyOrders() *buyOrders {
	head := newOrder(math.MinInt64, "", math.MaxFloat64, 0)
	return &buyOrders{head: head}
}

func (o *buyOrders) push(order *order) {
	curr, next := o.head, o.head.Next
	for next != nil && order.Price < next.Price {
		curr, next = next, next.Next
	}
	// now,
	// order.Price >= next.Price
	curr.Next, order.Next = order, next
}

func (o *buyOrders) pop() *order {
	if o.head.Next == nil {
		return nil
	}
	res := o.head.Next
	o.head.Next = o.head.Next.Next
	res.Next = nil
	return res
}

func (o *buyOrders) cancel(id int64) {
	curr := o.head
	for curr.Next != nil && curr.Next.ID != id {
		curr = curr.Next
	}
	if curr.Next == nil {
		return
	}
	curr.Next = curr.Next.Next
}

func (o *buyOrders) cancelAll() {
	o.head.Next = nil
}

func (o *buyOrders) canBuy(price float64) bool {
	return o.head.Next.Price >= price
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
		var orderBuf bytes.Buffer
		orderDec := gob.NewDecoder(&orderBuf)

		sellOrders, buyOrders := newSellOrders(), newBuyOrders()
		for {
			select {
			case <-exCtx.Done():
				log.Println("ctx done.")
			case msg := <-getOrder:
				// TODO: 把这里相关的部分，写成闭包
				orderBuf.Write(msg.Payload)
				var order order
				if err := orderDec.Decode(&order); err != nil {
					log.Fatal("decode order error:", err)
				}
				switch order.Type {
				case "sell":
					sellOrders.push(&order)
				case "buy":
					buyOrders.push(&order)
				default:
					log.Fatal("出现了错误的 order type")
				}
			}
		}
	}()

	return ex
}
