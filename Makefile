.PHONY:

build:
	go build -o ./.bin/bot cmd/bot/main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

run: build fmt vet
	./.bin/bot

cfg:
	go run scripts/createPrivateConfig.go

b: build
r: run
