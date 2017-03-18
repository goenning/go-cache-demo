package redis_test

import (
	"testing"

	"time"

	"github.com/goenning/go-cache-demo/cache/redis"
)

var redisURL = "redis://localhost:6379"

func parse(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func TestWrongURL(t *testing.T) {
	storage, err := redis.NewStorage("wrong://wtf")
	if err == nil || storage != nil {
		t.Fail()
	}
}

func TestGetEmpty(t *testing.T) {
	storage, _ := redis.NewStorage(redisURL)
	content := storage.Get("MY_KEY")

	assertContentEquals(t, content, "")
}

func TestGetValue(t *testing.T) {
	storage, _ := redis.NewStorage(redisURL)
	storage.Set("MY_KEY", "123456", parse("5s"))
	content := storage.Get("MY_KEY")

	assertContentEquals(t, content, "123456")
}

func TestGetExpiredValue(t *testing.T) {
	storage, _ := redis.NewStorage(redisURL)
	storage.Set("MY_KEY", "123456", parse("1s"))
	time.Sleep(parse("1s"))
	content := storage.Get("MY_KEY")

	assertContentEquals(t, content, "")
}

func assertContentEquals(t *testing.T, content, expected string) {
	if content != expected {
		t.Errorf("content should '%s', but was '%s'", expected, content)
	}
}
