#!/bin/bash

killall dvembed
git pull
rm -r downloads
#go mod sync
go test -v
go install dvembed
# TODO: add flag to continue and run
# if [ $1 ] ... fi
nohup dvembed &
echo "Is dvembed running?"
echo "Search for pid"
pidof dvembed
# ps aux | awk -F ' ' '{print $1}'
