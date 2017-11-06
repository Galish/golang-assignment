package indexer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Galish/golang-assignment/crawler"
	"github.com/Galish/golang-assignment/frontend"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

type Broker struct {
	instance broker.Broker
}

type Data struct {
	Result []crawler.Message
	Query  frontend.SearchQuery
}

var ch = make(chan Data)

func (b *Broker) init() {
	amqpBroker := rabbitmq.NewBroker(
		broker.Addrs(amqpAddr),
	)

	if err := amqpBroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := amqpBroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	fmt.Println("Broker listening:", amqpAddr)

	b.instance = amqpBroker
}

func (b *Broker) sub() {
	var err error

	err = b.subCrawler()
	if err != nil {
		fmt.Println(err)
	}

	err = b.subSearch()
	if err != nil {
		fmt.Println(err)
	}
}

func (b *Broker) subCrawler() error {
	_, err := b.instance.Subscribe(topicCrawler, func(p broker.Publication) error {
		// TODO: add header
		message := crawler.Message{}
		json.Unmarshal(p.Message().Body, &message)

		fmt.Printf("[sub] received message #%d\n", message.ID)

		key := getKey(message.ID)

		put(message.ID, key, p.Message().Body)
		index(key, []string{message.Title, message.HTML})

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *Broker) subSearch() error {
	_, err := b.instance.Subscribe(topicSearch, func(p broker.Publication) error {
		if p.Message().Header["Action"] != "search" {
			return nil
		}

		query := frontend.SearchQuery{}

		json.Unmarshal(p.Message().Body, &query)

		fmt.Printf("[sub] received search term #%s \"%s\"\n", query.ID, query.Term)

		result, err := find(query.Search)

		if err != nil {
			fmt.Println("search error:", err)
			return err
		}

		// fmt.Println("search done:", len(result))

		data := Data{
			Result: result,
			Query:  query,
		}

		ch <- data

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *Broker) pub() {
	for {
		data := <-ch
		query := data.Query

		result := frontend.SearchResult{
			ID:     query.ID,
			Term:   query.Term,
			Result: data.Result,
		}
		resultJSON, _ := json.Marshal(result)

		msg := &broker.Message{
			Header: map[string]string{
				"Action": "result",
			},
			Body: resultJSON,
		}

		if err := b.instance.Publish(topicSearch, msg); err != nil {
			log.Printf("[pub] failed: %v", err)
		} else {
			fmt.Printf("[pub] pubbed search result #%s (%d) \n", query.ID, len(data.Result))
		}
	}
}
