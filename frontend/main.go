package frontend

import (
	"github.com/Galish/golang-assignment/crawler"
	"github.com/micro/go-micro/broker"
)

var amqpBroker broker.Broker

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
	runBroker()

	sCh := make(chan Query)
	rCh := make(chan Search)

	go pub(sCh)
	go sub(rCh)

	serveWS(sCh, rCh)
}
