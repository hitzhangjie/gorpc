#!/bin/bash
CURR_PATH=$(cd "$(dirname "$0")"; pwd)

ulimit -c 1024000 -S

count=`ps -fe | grep {{.ServerName}} | grep -v grep | wc -l`
if [ $count -eq 0 ]
then
    if [ ! -f $CURR_PATH/{{.ServerName}} ]
    then
        echo "$CURR_PATH/{{.ServerName}} not exist"
        exit
    fi
    nohup  $CURR_PATH/{{.ServerName}} >> $CURR_PATH/../log/nohup.log 2>&1 &
    if [ $? -eq 0 ]
    then
        echo "{{.ServerName}} started"
    else
        echo "{{.ServerName}} started, check nohup.log"
    fi
else
    echo "{{.ServerName}} already runing"
fi
