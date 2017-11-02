package indexer

import "fmt"

var keyVal Rkv

func Run() {
	broker := Broker{}
	keyVal = Rkv{}

	keyVal.init()

	broker.init()
	go broker.sub()
	go broker.pub()

	fmt.Println("Indexer service is running")

	select {}
}
