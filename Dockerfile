FROM golang:1.24.6 AS builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o e ./cmd/e

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/e /usr/local/bin/e
CMD ["/usr/local/bin/e", "serve"]
