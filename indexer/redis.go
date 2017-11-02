package indexer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Galish/golang-assignment/crawler"
	"github.com/Galish/golang-assignment/frontend"
	"github.com/go-redis/redis"
)

type Rkv struct {
	Client *redis.Client
}

func (r *Rkv) init() {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (r *Rkv) Put(key string, value []byte) error {
	return r.Client.Set(key, value, 0).Err()
}

func (r *Rkv) Add(key string, value string) error {
	return r.Client.SAdd(key, value).Err()
}

func (r *Rkv) Get(key string) ([]byte, error) {
	return r.Client.Get(key).Bytes()
}

func (r *Rkv) GetKeys(key string) *redis.StringSliceCmd {
	return r.Client.SMembers(key)
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

func find(search map[string][]frontend.SearchTerm) ([]crawler.Message, error) {
	messages := []crawler.Message{}
	termIDs := make(map[string][]string)

	for key := range search {
		for i := range search[key] {
			term := search[key][i].Term
			ids := keyVal.GetKeys(term)
			err := ids.Err()

			if err != nil {
				fmt.Println(err)
				break
			}

			termIDs[term] = ids.Val()
		}
	}

	ids := getIDs(termIDs, "and")

	for _, id := range ids {
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

func getIDs(termIDs map[string][]string, exp string) []string {
	var slices [][]string

	for term := range termIDs {
		slices = append(slices, termIDs[term])
	}

	switch exp {
	case "or":
		return merge(slices)
	case "and":
		return inter(slices)
	}

	return []string{}
}
