package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	redis "gopkg.in/redis.v5"
)

var cached string
var when int64

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func index(w http.ResponseWriter, r *http.Request) {
	content := fmt.Sprintf("Hello World! Current time is: %s", time.Now())
	cached = content
	when = time.Now().Unix()

	w.Write([]byte(content))
}

func cache(seconds int64, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Unix()
		if when > 0 && (when+seconds >= now) {
			w.Write([]byte(cached))
		} else {
			handler(w, r)
		}
	})
}

func main() {
	h := cache(5, index)
	http.Handle("/", h)
	fmt.Println("Server is up and listening.")
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
