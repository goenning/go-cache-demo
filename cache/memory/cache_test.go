package memory_test

import (
	"testing"

	"time"

	"bytes"

	"github.com/goenning/go-cache-demo/cache/memory"
)

func parse(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func TestGetEmpty(t *testing.T) {
	storage := memory.NewStorage()
	content := storage.Get("MY_KEY")

	assertContentEquals(t, content, []byte(""))
}

func TestGetValue(t *testing.T) {
	storage := memory.NewStorage()
	storage.Set("MY_KEY", []byte("123456"), parse("5s"))
	content := storage.Get("MY_KEY")

	assertContentEquals(t, content, []byte("123456"))
}

func TestGetExpiredValue(t *testing.T) {
	storage := memory.NewStorage()
	storage.Set("MY_KEY", []byte("123456"), parse("1s"))
	time.Sleep(parse("1s200ms"))
	content := storage.Get("MY_KEY")

	assertContentEquals(t, content, []byte(""))
}

func assertContentEquals(t *testing.T, content, expected []byte) {
	if !bytes.Equal(content, expected) {
		t.Errorf("content should '%s', but was '%s'", expected, content)
	}
}
