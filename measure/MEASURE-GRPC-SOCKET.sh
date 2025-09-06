#!/bin/sh

ADDRESS=unix:@echo

echo-go serve grpc --socket @echo & 

./MEASURE.sh grpc grpc $ADDRESS

killall echo-go
