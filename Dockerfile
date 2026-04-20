FROM golang:1.26-alpine AS builder

WORKDIR /build

COPY go.mod .
COPY cmd/ ./cmd/

RUN go build -o /devops-reporter ./cmd/

FROM alpine:3.23

COPY --from=builder /devops-reporter /usr/local/bin/devops-reporter

ENTRYPOINT ["devops-reporter"]
