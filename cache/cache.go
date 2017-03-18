package cache

import "time"

//Storage mecanism for caching strings
type Storage interface {
	Get(key string) string
	Set(key, content string, duration time.Duration)
}
