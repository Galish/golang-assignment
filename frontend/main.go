package frontend

import "github.com/Galish/golang-assignment/crawler"

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

// Run Frontend service
func Run() {
	broker := Broker{}
	ws := WS{}

	chSearch := make(chan SearchQuery)
	chReslt := make(chan SearchResult)

	broker.init()
	go broker.pub(chSearch)
	go broker.sub(chReslt)

	ws.init(chSearch, chReslt)
	ws.serve()
}
