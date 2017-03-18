package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"strings"

	"github.com/goenning/go-cache-demo/cache"
	"github.com/goenning/go-cache-demo/cache/redis"
)

var storage cache.Storage

func init() {
	storage, _ = redis.NewStorage(os.Getenv("REDIS_URL"))
}

func index(w http.ResponseWriter, r *http.Request) {
	content := fmt.Sprintf("Hello World! Current time is: %s", time.Now())
	w.Write([]byte(content))
}

func isCacheable(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "text/html")
}

func cached(duration string, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !isCacheable(r) {
			handler(w, r)
			return
		}

		content := storage.Get(r.RequestURI)
		if content != "" {
			fmt.Print("Cache Hit!\n")
			w.Write([]byte(content))
		} else {
			c := httptest.NewRecorder()
			handler(c, r)

			for k, v := range c.HeaderMap {
				w.Header()[k] = v
			}

			w.WriteHeader(c.Code)
			content := c.Body.String()

			if d, err := time.ParseDuration(duration); err == nil {
				fmt.Printf("New cached page: %s for %s\n", r.RequestURI, duration)
				storage.Set(r.RequestURI, content, d)
			} else {
				fmt.Printf("Content not cached, %s\n", err)
			}

			w.Write([]byte(content))
		}

	})
}

func main() {
	http.Handle("/", cached("5s", index))
	fmt.Println("Server is up and listening.")
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
