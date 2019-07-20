#!/bin/bash
CURR_PATH=$(cd "$(dirname "$0")"; pwd)
echo $CURR_PATH

ulimit -c 1024000 -S

ps -fe|grep test_nrpc |grep -v grep
if [ $? -eq 0 ]
then
    nohup  $CURR_PATH/test_nrpc >> $CURR_PATH/../log/nohup.log 2>&1 &
else
    echo "test_nrpc runing....."
fi
