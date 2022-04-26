#!/bin/sh
nohup ./captcha_micro -conf prod -port 8099 > captcha.log &
nohup ./captcha_micro -conf prod -port 8098 > captcha.log &
nohup ./user_micro -conf prod -port 8096 > user.log &
nohup ./article_micro -conf prod -port 8097 > article.log &
tail -f user.log
