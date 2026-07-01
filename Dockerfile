FROM golang:1.26.4 AS builder
WORKDIR /app
COPY . ./
RUN apt-get update
RUN apt-get install unzip
RUN ./tools/fetch-protoc.sh
ENV PATH="/root/local/bin:${PATH}"
RUN make rpc
RUN make grpc
RUN make connect
RUN CGO_ENABLED=0 GOOS=linux go build -v -o echo-go .

FROM gcr.io/distroless/static-debian13
COPY --from=builder /app/echo-go /usr/local/bin/echo-go
CMD ["/usr/local/bin/echo-go", "serve",  "connect"]
