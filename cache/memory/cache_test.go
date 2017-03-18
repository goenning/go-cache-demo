package memory_test

import (
	"testing"

	"time"

	"github.com/goenning/go-cache-demo/cache/memory"
)

func parse(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func TestGetEmpty(t *testing.T) {
	storage := memory.NewStorage()
	content := storage.Get("MY_KEY")

	assertContentEquals(t, content, "")
}

func TestGetValue(t *testing.T) {
	storage := memory.NewStorage()
	storage.Set("MY_KEY", "123456", parse("5s"))
	content := storage.Get("MY_KEY")

	assertContentEquals(t, content, "123456")
}

func TestGetExpiredValue(t *testing.T) {
	storage := memory.NewStorage()
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
