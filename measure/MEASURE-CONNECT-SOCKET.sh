#!/bin/sh

ADDRESS=unix:@echoconnect

echo-go serve connect --socket @echoconnect & 

./MEASURE.sh connect grpc $ADDRESS
./MEASURE.sh connect connect $ADDRESS
./MEASURE.sh connect connect-grpc $ADDRESS
./MEASURE.sh connect connect-grpc-web $ADDRESS

killall echo-go
