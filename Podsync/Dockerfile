FROM golang:alpine3.10

WORKDIR /apps
COPY go.mod go.sum ./

RUN wget -O /usr/bin/youtube-dl https://github.com/ytdl-org/youtube-dl/releases/download/2019.11.05/youtube-dl && \
    chmod +x /usr/bin/youtube-dl && \
    apk --no-cache add ca-certificates python ffmpeg
RUN go mod download

COPY ./cmd ./
COPY ./docs .
COPY ./pkg .

RUN go build -o main .

EXPOSE 8080

CMD ["./main", "--config config.toml"]
