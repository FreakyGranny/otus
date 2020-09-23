FROM debian:buster-slim

RUN mkdir -p /app
WORKDIR /app

ENV WAIT_VERSION 2.7.3
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

COPY bin/calendar /app
COPY configs/calendar_config.yaml /app/config.yaml

CMD ["/app/calendar", "-h"]
