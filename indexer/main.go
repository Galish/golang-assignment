package indexer

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Galish/golang-assignment/crawler"
	"github.com/Galish/golang-assignment/frontend"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-plugins/broker/rabbitmq"
)

const (
	amqpAddr     = "amqp://localhost"
	topicCrawler = "topic.crawler"
	topicSearch  = "topic.search"
)

var (
	amqpBroker broker.Broker
	keyVal     Rkv
)

func sub() {
	fmt.Println("Broker listening:", amqpAddr)
	var subErr error

	_, subErr = amqpBroker.Subscribe(topicCrawler, func(p broker.Publication) error {
		message := crawler.Message{}
		json.Unmarshal(p.Message().Body, &message)

		fmt.Printf("[sub] received message #%d\n", message.ID)

		key := getKey(message.ID)

		put(message.ID, key, p.Message().Body)
		index(key, message.HTML)

		return nil
	})

	if subErr != nil {
		fmt.Println(subErr)
	}

	_, subErr = amqpBroker.Subscribe(topicSearch, func(p broker.Publication) error {
		search := frontend.Search{}
		json.Unmarshal(p.Message().Body, &search)
		id := p.Message().Header["ID"]
		term := search.Term
		res := search.Result

		if res != nil {
			return nil
		}

		fmt.Printf("[sub] received search term #%s \"%s\"\n", id, term)

		messages, err := find(search.Term)

		if err != nil {
			fmt.Println("search error:", err)
			return err
		}

		result := frontend.Search{
			Term:   term,
			Result: messages,
		}
		resultJSON, _ := json.Marshal(result)

		msg := &broker.Message{
			Header: map[string]string{
				"ID": id,
			},
			Body: resultJSON,
		}

		if err := amqpBroker.Publish(topicSearch, msg); err != nil {
			log.Printf("[pub] failed: %v", err)
		} else {
			fmt.Printf("[pub] pubbed search result #%s (%d) \n", id, len(messages))
		}

		return nil
	})

	if subErr != nil {
		fmt.Println(subErr)
	}
}

func put(id int, key string, value []byte) {
	err := keyVal.Put(key, value)

	if err != nil {
		fmt.Println("put err", err)
	} else {
		fmt.Printf("[put] #%d\n", id)
	}
}

func index(key string, HTML string) {
	tokens := parseTokens(strings.NewReader(HTML))

	for _, token := range tokens {
		err := keyVal.Add(token, key)

		if err != nil {
			fmt.Println("index err", err)
		}
		// else {
		// 	fmt.Printf("[index] %s\n", token)
		// }
	}
}

func find(term string) ([]crawler.Message, error) {
	// var messages []crawler.Message
	// keyVal := redis.NewKV()
	messages := []crawler.Message{}

	ids := keyVal.GetKeys(term)
	err := ids.Err()

	if err != nil {
		return nil, err
	}

	for _, id := range ids.Val() {
		message := crawler.Message{}
		val, err := keyVal.Get(id)

		if err != nil {
			return nil, err
		}

		json.Unmarshal(val, &message)

		messages = append(messages, message)
	}

	return messages, nil
}

// Run Indexer service
func Run() {
	// Initiate RabbitMQ broker
	amqpBroker = rabbitmq.NewBroker(
		broker.Addrs(amqpAddr),
	)

	if err := amqpBroker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := amqpBroker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	// Initiate Redis key-value storage
	keyVal = NewKV(
		"localhost:6379",
		"",
		0,
	)

	fmt.Println("Indexer service is running")

	sub()

	select {}
}
