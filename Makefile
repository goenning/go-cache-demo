test: 
	go test $$(go list ./... | grep -v /vendor/) -cover

run: 
	PORT=8080 REDIS_URL=redis://localhost:6379 go run main.go