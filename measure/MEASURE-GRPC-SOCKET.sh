#!/bin/sh

ADDRESS=unix:@echogrpc

echo-go serve grpc --socket @echogrpc & 

./MEASURE.sh grpc grpc $ADDRESS
./MEASURE.sh grpc connect-grpc $ADDRESS

killall echo-go
