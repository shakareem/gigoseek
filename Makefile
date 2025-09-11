.PHONY:

build:
	go build -o ./.bin/bot cmd/bot/main.go

b: build

fmt:
	go fmt ./...

vet:
	go vet ./...

run: build fmt vet
	./.bin/bot

r: run
