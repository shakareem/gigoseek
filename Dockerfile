FROM golang:1.24.6

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
COPY pkg pkg
COPY configs configs

RUN go build -v -o ./.bin/bot ./cmd/bot

CMD ["./.bin/bot"]
