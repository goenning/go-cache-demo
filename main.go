package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

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
	// Instead of sleeping, imagine that this handler takes 2 seconds to query database and build the result...
	time.Sleep(2 * time.Second)

	content := fmt.Sprintf(`
		<h1>Hello World! You are on: %s</h1> 
		<p>Current time is: %s</p>

		<ul>
			<li><a href="/">Home</a></li>
			<li><a href="/?page=1">Page 1</a></li>
			<li><a href="/?page=2">Page 2</a></li>
			<li><a href="/about">About (not cached!)</a></li>
		</ul>
	`, r.RequestURI, time.Now())
	w.Write([]byte(content))
}

func about(w http.ResponseWriter, r *http.Request) {
	// Instead of sleeping, imagine that this handler takes 2 seconds to query database and build the result...
	time.Sleep(2 * time.Second)

	content := fmt.Sprintf(`
		<h1>About!</h1> 
		<p>Current time is: %s</p>

		<ul>
			<li><a href="/">Home</a></li>
			<li><a href="/about">About (not cached!)</a></li>
		</ul>
	`, time.Now())
	w.Write([]byte(content))
}

func main() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	http.Handle("/", cached("10s", index))
	http.HandleFunc("/about", about)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is up and listening on port %s.\n", port)
	http.ListenAndServe(":"+port, nil)
}
