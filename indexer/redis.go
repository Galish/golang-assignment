package indexer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Galish/golang-assignment/crawler"
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

func index(key string, texts []string) {
	for _, text := range texts {
		tokens := parseTokens(strings.NewReader(text))

		for _, token := range tokens {
			err := keyVal.Add(token, key)

			if err != nil {
				fmt.Println("index err", err)
			}
			//  else {
			// 	fmt.Printf("[index] %s - %s\n", token, key)
			// }
		}
	}
}

func find(search interface{}) ([]crawler.Message, error) {
	messages := []crawler.Message{}
	ids, _ := findIDs(search)

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

func findIDs(s interface{}) ([]string, []string) {
	m := s.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			if k == "term" {
				keys := keyVal.GetKeys(vv)
				err := keys.Err()

				if err != nil {
					fmt.Println(err)
					break
				}

				terms := keys.Val()

				return nil, terms
			}
		case []interface{}:
			terms := [][]string{}

			for _, i := range vv {
				_ids, _terms := findIDs(i)

				if _ids != nil {
					terms = append(terms, _ids)
				} else {
					terms = append(terms, _terms)
				}
			}

			ids := getIDs(terms, k)

			return ids, nil
		}
	}

	return nil, nil
}

func getIDs(slices [][]string, exp string) []string {
	switch exp {
	case "or":
		return merge(slices)
	case "and":
		return inter(slices)
	}

	return []string{}
}
