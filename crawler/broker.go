package crawler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

type Broker struct {
	instance broker.Broker
}

func (b *Broker) init() {
	amqpBroker := rabbitmq.NewBroker(
		broker.Addrs(amqpAddr),
	)

	if err := amqpBroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := amqpBroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	fmt.Println("Broker is running at:", amqpAddr)

	b.instance = amqpBroker
}

func (b *Broker) pub() {
	for {
		message := <-ch
		messageJSON, _ := json.Marshal(message)

		msg := &broker.Message{
			Header: map[string]string{
				"Index": fmt.Sprintf("%d", postIndex),
			},
			Body: messageJSON,
		}

		if err := b.instance.Publish(topic, msg); err != nil {
			log.Printf("[pub] failed: %v", err)
		} else {
			fmt.Printf("[pub] pubbed message #%d (index: %d)\n", message.ID, message.Index)
		}
	}
}
