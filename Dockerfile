FROM golang:1.25.1 AS builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o echo-go .

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/echo-go /usr/local/bin/echo-go
CMD ["/usr/local/bin/echo-go", "serve",  "connect"]
