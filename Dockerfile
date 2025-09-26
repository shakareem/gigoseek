FROM golang:1.24.6

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
COPY pkg pkg
COPY configs configs
COPY scripts scripts

RUN go build -o ./.bin/bot ./cmd/bot

CMD ["./.bin/bot"]
