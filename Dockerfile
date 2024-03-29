FROM golang:1.20 as builder

RUN mkdir /api
COPY . /api

WORKDIR /api
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bot-api ./main

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /api/bot-api /usr/bin/bot-api
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/bot-api"]
