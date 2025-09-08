#!/bin/sh

ADDRESS=unix:@echo

echo-go serve connect --socket @echo & 

./MEASURE.sh connect grpc $ADDRESS

killall echo-go
