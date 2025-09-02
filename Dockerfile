FROM golang:1.24.6 as builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o echo-server ./cmd/echo-server

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/echo-server /usr/local/bin/echo-server
COPY --from=builder /app/data /data
CMD ["/usr/local/bin/echo-server"]
