.PHONY: build run all
.DEFAULT_GOAL := all

build:
	go build -i -o build/server backend/main.go

run:
	./build/server

test:
	go test  ./...

all: build run