package frontend

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/micro/go-micro/broker"
)

func pub(ch <-chan Query) {
	for {
		// fmt.Println("PUB+++", <-srchCh)
		query := <-ch
		fmt.Println("!!!", query)

		id := "12345"
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

		if err := amqpBroker.Publish(topicSearch, msg); err != nil {
			log.Printf("[pub] failed: %v", err)
		} else {
			fmt.Printf("[pub] pubbed search term #%s \"%s\"\n", id, term)
		}
	}
}

func sub(ch chan<- Search) {
	fmt.Println("Broker listening:", amqpAddr)

	_, err := amqpBroker.Subscribe(topicSearch, func(p broker.Publication) error {
		search := Search{}
		json.Unmarshal(p.Message().Body, &search)
		id := p.Message().Header["ID"]
		term := search.Term
		res := search.Result

		if res == nil {
			return nil
		}

		log.Printf("[sub] received search results #%s \"%s\" (%d)\n", id, term, len(res))

		ch <- search

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
