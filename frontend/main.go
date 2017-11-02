package frontend

// Run Frontend service
func Run() {
	broker := Broker{}
	ws := WS{}

	broker.init()
	go broker.pub()
	go broker.sub()

	ws.init(broker)
	ws.serve()
}
