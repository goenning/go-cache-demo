package cache

import "time"

//Storage mecanism for caching strings
type Storage interface {
	Get(key string, duration time.Duration) string
	Set(key, content string, duration time.Duration)
}
