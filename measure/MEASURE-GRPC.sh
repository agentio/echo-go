#!/bin/sh

ADDRESS=localhost:8080

echo-go serve grpc --port 8080 &

./MEASURE.sh grpc grpc $ADDRESS
./MEASURE.sh grpc connect-grpc $ADDRESS

killall -9 echo-go
