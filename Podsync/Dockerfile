#!/bin/bash

FROM golang:alpine3.10

WORKDIR /apps
COPY go.mod go.sum ./

RUN wget -O /usr/bin/youtube-dl https://github.com/ytdl-org/youtube-dl/releases/download/2019.11.05/youtube-dl
RUN chmod +x /usr/bin/youtube-dl

RUN echo "ipv6" >> /etc/modules
RUN apk --no-cache add ca-certificates python ffmpeg git

RUN go mod download

COPY ./cmd/podsync ./cmd/podsync
COPY ./docs ./docs
COPY ./pkg ./pkg
COPY config.toml .

# RUN go build -o main .
# COPY main .

EXPOSE 8080

# CMD ["main", "--config config.toml"]
CMD ["go", "run", "./cmd/podsync"]
# CMD ["ls"]
