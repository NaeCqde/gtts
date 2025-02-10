FROM golang:1-alpine AS builder


WORKDIR /app

COPY ./dca /app/dca
COPY ./disgolf /app/disgolf
COPY *.go go.mod go.sum /app/

RUN go mod download
RUN go build -o gtts .


FROM alpine:latest


WORKDIR /app

RUN apk --no-cache add ffmpeg

COPY --from=builder /app/gtts .
RUN chmod +x gtts

ENTRYPOINT ./gtts
