#!/bin/sh

ADDRESS=mccarthy.lan:8080

./MEASURE.sh grpc grpc $ADDRESS
./MEASURE.sh grpc connect-grpc $ADDRESS

