package indexer

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Galish/golang-assignment/crawler"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

const (
	amqpAddr = "amqp://localhost"
	topic    = "topic.crawler"
)

var (
	amqpBroker broker.Broker
	keyVal     Rkv
)

func subscribe() {
	fmt.Println("Broker listening:", amqpAddr)

	_, err := amqpBroker.Subscribe(topic, func(p broker.Publication) error {
		message := crawler.Message{}
		json.Unmarshal(p.Message().Body, &message)

		fmt.Printf("[sub] received message #%d\n", message.ID)

		key := getKey(message.ID)

		put(key, p.Message().Body)
		index(key, message.HTML)

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func put(key string, value []byte) {
	err := keyVal.Put(key, value)

	if err != nil {
		fmt.Println("put err", err)
	} else {
		fmt.Printf("[put] #%s\n", key)
	}
}

func index(key string, HTML string) {
	tokens := parseTokens(strings.NewReader(HTML))

	for _, token := range tokens {
		err := keyVal.Add(token, key)

		if err != nil {
			fmt.Println("index err", err)
		}
		// else {
		// 	fmt.Printf("[index] %s\n", token)
		// }
	}
}

// Run Indexer service
func Run() {
	// Initiate RabbitMQ broker
	amqpBroker = rabbitmq.NewBroker(
		broker.Addrs(amqpAddr),
	)

	if err := amqpBroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := amqpBroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	// Initiate Redis key-value storage
	keyVal = NewKV(
		"localhost:6379",
		"",
		0,
	)

	fmt.Println("Indexer service is running")

	subscribe()

	select {}
}
