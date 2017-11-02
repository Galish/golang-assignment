package frontend

import (
	"github.com/Galish/golang-assignment/crawler"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/broker"
)

type Broker struct {
	instance broker.Broker
	chSearch chan Query
	chReslt  chan Search
}

type WS struct {
	conn     *websocket.Conn
	chSearch chan Query
	chReslt  chan Search
}

type Search struct {
	Term   string            `json:"term"`
	Result []crawler.Message `json:"result"`
}

type Query struct {
	Search string `json:"search"`
}
