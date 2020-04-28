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

// TODO: 添加 tick 的属性
type tick struct {
	ID     int64
	Price  float64
	Amount float64
}

type pubsub interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(topic string) (<-chan *message.Message, error)
	Close() error
}

type order struct {
	ID int64
	// TODO: 把 OrderType 的 string 转变成枚举
	Type          string
	Price, Amount float64
	Next          *order
}

func newOrder(id int64, orderType string, price, quantity float64) *order {
	return &order{
		ID:     id,
		Type:   orderType,
		Price:  price,
		Amount: quantity,
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
	return o.head.Next != nil &&
		o.head.Next.Price <= price
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
		// about order
		getOrder, err := pubsub.Subscribe("order")
		if err != nil {
			panic("pubsub.Subscribe(order) error" + err.Error())
		}
		var orderBuf bytes.Buffer
		orderDec := gob.NewDecoder(&orderBuf)
		sellOrders, buyOrders := newSellOrders(), newBuyOrders()
		// about tick
		getTick, err := pubsub.Subscribe("tick")
		if err != nil {
			panic("pubsub.Subscribe(tick) error" + err.Error())
		}
		var tickBuf bytes.Buffer
		tickDec := gob.NewDecoder(&tickBuf)
		for {
			select {
			case <-exCtx.Done():
				log.Println("ctx done.")
			case msg := <-getOrder:
				// TODO: 把 order 的处理方式写成闭包
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
			case msg := <-getTick:
				// TODO: 把 tick 的处理方式写成闭包
				tickBuf.Write(msg.Payload)
				var tick tick
				if err := tickDec.Decode(&tick); err != nil {
					log.Fatal("decode order error:", err)
				}
				for tick.Amount != 0 && sellOrders.canSell(tick.Price) {
					order := sellOrders.pop()
					var trade Trade
					// 此时 order.Price >= tick.Price
					// 选用 tick.Price 的原因是为了得到更悲观的结果。
					// 因为没有盘口信息，所以无法根据订单的类型选择正确的成交价格。
					trade.Price = tick.Price
					if order.Amount > tick.Amount {
						trade.Amount = tick.Amount
						// 把没有执行完的订单，放回去
						order.Amount -= tick.Amount
						sellOrders.push(order)
					} else {
						trade.Amount = order.Amount
					}
					// 更新 tick 的信息，以便下一轮判断是否能够执行。
					tick.Amount -= trade.Amount
					// TODO: 更新交易所内的个人帐户信息
					// TODO: 发布成交信息
				}
				for tick.Amount != 0 && buyOrders.canBuy(tick.Price) {
					order := buyOrders.pop()
					var trade Trade
					// 此时 order.Price <= tick.Price
					// 选用 tick.Price 的原因是为了得到更悲观的结果。
					// 因为没有盘口信息，所以无法根据订单的类型选择正确的成交价格。
					trade.Price = order.Price
					if order.Amount > tick.Amount {
						trade.Amount = tick.Amount
						// 把没有执行完的订单，放回去
						order.Amount -= tick.Amount
						buyOrders.push(order)
					} else {
						trade.Amount = order.Amount
					}
					// 更新 tick 的信息，以便下一轮判断是否能够执行。
					tick.Amount -= trade.Amount
					// TODO: 更新交易所内的个人帐户信息
					// TODO: 发布成交信息
				}
			}
		}
	}()

	return ex
}

// Trade 是交易记录
type Trade struct {
	ID            int64
	Price, Amount float64
}

// Run TODO: 从数据库读取数据，并逐条发送到 pubsub
func (ex *exchange) Run(ctx context.Context, source <-chan tick, pubsub pubsub) {

}
