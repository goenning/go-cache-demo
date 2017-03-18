package redis

import (
	"time"

	r "gopkg.in/redis.v5"
)

var preffix = "_PAGE_CACHE_"

//Storage mecanism for caching strings in memory
type Storage struct {
	client *r.Client
}

//NewStorage creates a new redis storage
func NewStorage(url string) (*Storage, error) {
	var (
		opts *r.Options
		err  error
	)

	if opts, err = r.ParseURL(url); err != nil {
		return nil, err
	}

	return &Storage{
		client: r.NewClient(opts),
	}, nil
}

//Get a cached content by key
func (s Storage) Get(key string) string {
	val, _ := s.client.Get(preffix + key).Result()
	return val
}

//Set a cached content by key
func (s Storage) Set(key, content string, duration time.Duration) {
	s.client.Set(preffix+key, content, duration)
}
