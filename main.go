package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"strings"

	"github.com/goenning/go-cache-demo/cache"
	"github.com/goenning/go-cache-demo/cache/memory"
	"github.com/goenning/go-cache-demo/cache/redis"
)

var storage cache.Storage

func init() {
	strategy := flag.String("s", "memory", "Cache strategy (memory or redis)")
	flag.Parse()

	if *strategy == "memory" {
		storage = memory.NewStorage()
	} else if *strategy == "redis" {
		var err error
		if storage, err = redis.NewStorage(os.Getenv("REDIS_URL")); err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Sprintf("Invalid cache strategy %s.", *strategy))
	}
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is up and listening on port %s.\n", port)
	http.ListenAndServe(":"+port, nil)
}
