#!/bin/sh
nohup ./gateway_micro -conf prod > main.log &
tail -F main.log
