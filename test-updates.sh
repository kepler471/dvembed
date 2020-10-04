#!/bin/bash

killall dvembed
git pull
rm -r downloads
#go mod sync
go test -v
go run dvembed
# TODO: a different script to install and run, probably includes the following lines
# if [ $1 ] ... fi
# nohup dvembed &
# echo "Is dvembed running?"
# echo "Search for pid"
# pidof dvembed
# ps aux | awk -F ' ' '{print $1}'
