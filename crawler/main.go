package crawler

import "fmt"

type Message struct {
	ID     int    `json:"id"`
	Index  int    `json:"index"`
	Link   string `json:"link"`
	Avatar string `json:"avatar"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	HTML   string `json:"html"`
	Text   string `json:"text"`
}

var ch = make(chan Message)

// Run Crawler service
func Run() {
	fmt.Println("Crawler service is running")

	broker := Broker{}
	crawler := Crawler{}

	broker.init()
	go broker.pub()

	crawler.init()
	crawler.run()

	// select {}
}
