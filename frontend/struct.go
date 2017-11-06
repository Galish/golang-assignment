package frontend

import (
	"github.com/Galish/golang-assignment/crawler"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/broker"
)

type Broker struct {
	instance broker.Broker
}

type WS struct {
	conn *websocket.Conn
}

type SearchResult struct {
	ID     string            `json:"id"`
	Term   string            `json:"term"`
	Result []crawler.Message `json:"result"`
}

type SearchQuery struct {
	ID     string      `json:"id"`
	Term   string      `json:"term"`
	Search interface{} `json:"search"`
}

type SearchTerm struct {
	Term string `json:"term"`
}
