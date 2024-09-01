.PHONY: run 
run:
	go run main.go

.PHONY: build
build:
	go build -o bin/main main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: bench
bench:
	go test -bench=. ./... -benchmem

.PHONY: clean
clean:
	rm -rf bin

.PHONY: docker-build
docker-build:
	docker build -t abc-trading .
