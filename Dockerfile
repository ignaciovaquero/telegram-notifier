FROM golang:1.14.1-alpine3.11 AS builder

WORKDIR /go/src/github.com/igvaquero18/telegram-notifier
COPY . .
RUN go build


FROM alpine:3.11.3

LABEL org.opencontainers.image.authors="Ignacio Vaquero <ivaqueroguisasola@gmail.com>" \
      org.opencontainers.image.source="https://github.com/igvaquero18/telegram-notifier" \
      org.opencontainers.image.title="Telegram Notifier Bot" \
      org.opencontainers.image.description="Image for the Telegram Notifier bot" \
      org.opencontainers.image.version="0.1.0"

ENV NOTIFIER_BOT_LISTEN_ADDRESS="0.0.0.0" \
    NOTIFIER_BOT_LISTEN_PORT=8081 \
    NOTIFIER_BOT_TOKEN="" \
    NOTIFIER_BOT_TIMEOUT="15s" \
    NOTIFIER_BOT_VERBOSE="false"

COPY --from=builder /go/src/github.com/igvaquero18/telegram-notifier/telegram-notifier /go/bin/telegram-notifier

RUN apk update && \
    apk add tini=0.18.0-r0

EXPOSE 8081/tcp

ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/go/bin/telegram-notifier", "start"]
