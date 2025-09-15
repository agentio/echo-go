#!/bin/sh

SERVE_STACK=$1
CALL_STACK=$2
ADDRESS=$3

N=10000

tg=`echo-go call get --address=$ADDRESS -n $N --stack $CALL_STACK`
te=`echo-go call expand --address=$ADDRESS -n $N --stack $CALL_STACK`
tc=`echo-go call collect --address=$ADDRESS -n $N --stack $CALL_STACK`
ts=`echo-go call update --address=$ADDRESS -n $N --stack $CALL_STACK`

echo "| $SERVE_STACK | $CALL_STACK | $ADDRESS | $tg | $te | $tc | $ts |"
