#!/bin/sh
#nohup ./captcha_micro -conf prod -port 8089 > micros.log &
nohup ./captcha_micro -conf prod -port 8088 > micros.log &
nohup ./user_micro -conf prod -port 8086 > micros.log &
nohup ./article_micro -conf prod -port 8087 > micros.log &
tail -F micros.log
