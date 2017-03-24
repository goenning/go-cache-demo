package cache

import "time"

//Storage mecanism for caching strings
type Storage interface {
	Get(key string) []byte
	Set(key string, content []byte, duration time.Duration)
}
