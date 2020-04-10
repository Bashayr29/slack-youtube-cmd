FROM golang:1.12.7 AS builder

RUN mkdir /slack-youtube-cmd

COPY . /slack-youtube-cmd

WORKDIR /slack-youtube-cmd/main

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /slack-youtube-cmd/main

COPY --from=builder /slack-youtube-cmd/main .

EXPOSE 8080

CMD ["./main"]