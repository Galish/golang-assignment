package crawler

import "fmt"

// Run Crawler service
func Run() {
	fmt.Println("Crawler service is running")

	var messageChan = make(chan Message)

	broker := NewBroker(amqpAddr, messageChan)
	broker.Run()

	crawler := Crawler{}

	crawler.init()
	crawler.run()
}
