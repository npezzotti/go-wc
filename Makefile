NAME=go-wc
.DEFAULT_GOAL := build

build: clean fmt
	go build -o bin/$(NAME) cmd/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

fmt:
	go fmt ./...
