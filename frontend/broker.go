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
}

func (b *Broker) pub(ch <-chan SearchQuery) {
	for {
		query := <-ch

		queryJSON, err := json.Marshal(query)

		if err != nil {
			fmt.Println(err)
		}

		msg := &broker.Message{
			Header: map[string]string{
				"Action": "search",
			},
			Body: queryJSON,
		}

		fmt.Printf("!!!!!%+v\n", query)

		if err := b.instance.Publish(topicSearch, msg); err != nil {
			fmt.Printf("[pub] failed: %v", err)
		} else {
			// fmt.Printf("[pub] pubbed search term #%s \"%s\"\n", id, term)
			fmt.Printf("[pub] pubbed search term \"%s\"\n", query.Term)
		}
	}
}

func (b *Broker) sub(ch chan<- SearchResult) {
	fmt.Println("> Broker listening:", amqpAddr)

	_, err := b.instance.Subscribe(topicSearch, func(p broker.Publication) error {
		if p.Message().Header["Action"] != "result" {
			return nil
		}

		res := SearchResult{}

		json.Unmarshal(p.Message().Body, &res)

		if res.Result == nil {
			return nil
		}

		// fmt.Printf("[sub] received search results #%s \"%s\" (%d)\n", id, term, len(result))
		fmt.Printf("[sub] received search results \"%s\" (%d)\n", res.Term, len(res.Result))

		ch <- res

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
