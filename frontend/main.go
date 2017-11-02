package frontend

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
