# Run locally
.PHONY: run 
run:
	go run main.go

# Build the binary locally
.PHONY: build
build:
	go build -o bin/abc-trading main.go

# Play golang tests
.PHONY: test
test:
	go test -v ./...

# Clean binaries
.PHONY: clean
clean:
	rm -rf bin

# Run the docker container
.PHONY: docker-build
docker-build:
	docker build -t abc-trading .
