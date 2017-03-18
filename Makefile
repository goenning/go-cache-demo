test: 
	go test $$(go list ./... | grep -v /vendor/) -cover

run-redis: 
	PORT=8080 REDIS_URL=redis://localhost:6379 go run main.go -s redis

run-memory: 
	PORT=8080 go run main.go -s memory