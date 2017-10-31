package redis

import (
	"github.com/go-redis/redis"
)

type Rkv struct {
	Client *redis.Client
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

func NewKV(addr string, pass string, db int) Rkv {
	return Rkv{
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       db,
		}),
	}
}
