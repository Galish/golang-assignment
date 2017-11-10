package crawler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

type Broker struct {
	AmqpAddr   string
	AmqpBroker broker.Broker
	Messages   chan Message
}

func NewBroker(amqpAddr string, messageChan chan Message) *Broker {
	amqpBroker := rabbitmq.NewBroker(
		broker.Addrs(amqpAddr),
	)
	return &Broker{
		AmqpAddr:   amqpAddr,
		AmqpBroker: amqpBroker,
		Messages:   messageChan,
	}
}

func (b *Broker) Run() {
	b.init()
	go b.pub()
}

func (b *Broker) init() {
	if err := b.AmqpBroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := b.AmqpBroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	fmt.Println("Broker is running at:", b.AmqpAddr)
}

func (b *Broker) pub() {
	for {
		message := <-b.Messages
		messageJSON, _ := json.Marshal(message)

		msg := &broker.Message{
			Header: map[string]string{
				"Index": fmt.Sprintf("%d", postIndex),
			},
			Body: messageJSON,
		}

		if err := b.AmqpBroker.Publish(topic, msg); err != nil {
			log.Printf("[pub] failed: %v", err)
		} else {
			fmt.Printf("[pub] pubbed message #%d (index: %d)\n", message.ID, message.Index)
		}
	}
}
