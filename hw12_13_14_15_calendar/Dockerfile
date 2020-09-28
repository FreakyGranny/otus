FROM golang:1.14-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app
COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/calendar ./cmd/calendar

ENV WAIT_VERSION 2.7.3
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

FROM alpine:3.12

WORKDIR /app

COPY --from=builder /wait /wait
COPY --from=builder /go/bin/calendar /app/calendar
COPY configs/calendar_config.yaml /app/config.yaml

CMD ["/app/calendar", "-h"]
