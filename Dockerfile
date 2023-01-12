FROM golang:alpine AS builder

COPY . /throttler/
WORKDIR /throttler/

RUN go mod download
RUN go build -o ./bin/service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /throttler/bin/service .
COPY --from=builder /throttler/config/config.json ./config/
COPY --from=builder /throttler/docs ./docs

CMD ["./service", "run"]