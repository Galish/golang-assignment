package main

import (
	"flag"
	"fmt"

	"github.com/Galish/golang-assignment/crawler"
	"github.com/Galish/golang-assignment/frontend"
	"github.com/Galish/golang-assignment/indexer"
)

func main() {
	service := flag.String("service", "", "a string")
	flag.Parse()

	switch *service {
	case "indexer":
		indexer.Run()
	case "crawler":
		crawler.Run()
	case "frontend":
		frontend.Run()
	default:
		fmt.Println("Please specify a service to run:")
		fmt.Println("$ go run main.go -service=indexer/crawler/frontend")
	}
}
