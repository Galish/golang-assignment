package frontend

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Galish/golang-assignment/crawler"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
	// "github.com/micro/go-plugins/transport/rabbitmq"
)

// var amqpTransport transport.Transport
// var amqpClient transport.Client
var amqpBroker broker.Broker

const (
	amqpAddr    = "amqp://localhost"
	topicSearch = "topic.search"
)

type Search struct {
	Term   string
	Result []crawler.Message
}

func pub() {
	id := "12345"
	term := "all"

	search := Search{
		Term:   term,
		Result: nil,
	}
	searchJSON, _ := json.Marshal(search)

	msg := &broker.Message{
		Header: map[string]string{
			"ID": id,
		},
		Body: searchJSON,
	}

	if err := amqpBroker.Publish(topicSearch, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	} else {
		fmt.Printf("[pub] pubbed search term #%s \"%s\"\n", id, term)
	}
}

func sub() {
	_, err := amqpBroker.Subscribe(topicSearch, func(p broker.Publication) error {
		search := Search{}
		json.Unmarshal(p.Message().Body, &search)
		id := p.Message().Header["ID"]
		term := search.Term
		res := search.Result

		if res == nil {
			return nil
		}

		log.Printf("[sub] received search results #%s \"%s\" (%d)\n", id, term, len(res))

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

// Run Frontend service
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

	fmt.Println("Frontend service is running")

	sub()
	pub()

	select {}
}
