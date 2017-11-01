package frontend

import (
	"fmt"
	"log"

	"github.com/Galish/golang-assignment/crawler"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
	// "github.com/micro/go-plugins/transport/rabbitmq"
)

// var amqpTransport transport.Transport
// var amqpClient transport.Client
var (
	amqpBroker broker.Broker
	// srchCh     = make(chan Query)
	// rsltCh     = make(chan Search)
)

const (
	amqpAddr    = "amqp://localhost"
	topicSearch = "topic.search"
)

type Search struct {
	Term   string            `json:"term"`
	Result []crawler.Message `json:"result"`
}

type Query struct {
	Search string `json:"search"`
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

	sCh := make(chan Query)
	rCh := make(chan Search)

	go serveWS(sCh, rCh)
	go pub(sCh)
	go sub(rCh)

	select {}
}
