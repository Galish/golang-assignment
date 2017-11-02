package frontend

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

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

	fmt.Println("> Frontend service is running")

	b.instance = amqpBroker
	b.chSearch = make(chan Query)
	b.chReslt = make(chan Search)
}

func (b *Broker) pub() {
	for {
		query := <-b.chSearch
		id := "12345" // TODO: add ID to search query
		term := query.Search

		search := Search{
			Term:   term,
			Result: nil,
		}

		searchJSON, err := json.Marshal(search)

		if err != nil {
			fmt.Println(err)
		}

		msg := &broker.Message{
			Header: map[string]string{
				"ID": id,
			},
			Body: searchJSON,
		}

		if err := b.instance.Publish(topicSearch, msg); err != nil {
			fmt.Printf("[pub] failed: %v", err)
		} else {
			fmt.Printf("[pub] pubbed search term #%s \"%s\"\n", id, term)
		}
	}
}

func (b *Broker) sub() {
	fmt.Println("> Broker listening:", amqpAddr)

	_, err := b.instance.Subscribe(topicSearch, func(p broker.Publication) error {
		search := Search{}
		json.Unmarshal(p.Message().Body, &search)
		id := p.Message().Header["ID"]
		term := search.Term
		res := search.Result

		if res == nil {
			return nil
		}

		fmt.Printf("[sub] received search results #%s \"%s\" (%d)\n", id, term, len(res))

		b.chReslt <- search

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
