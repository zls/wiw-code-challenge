.PHONY: build run all
.DEFAULT_GOAL := all

build:
	go build -i -o build/server backend/main.go

run:
	./build/server

all: build run