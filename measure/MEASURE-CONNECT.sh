#!/bin/sh

ADDRESS=localhost:8080

echo-go serve connect --port 8080 &

./MEASURE.sh connect grpc $ADDRESS
./MEASURE.sh connect connect $ADDRESS
./MEASURE.sh connect connect-grpc $ADDRESS
./MEASURE.sh connect connect-grpc-web $ADDRESS

killall -9 echo-go
