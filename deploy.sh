#!/bin/bash

USER="root"
SERVER="10.33.106.127"
FILE="webcron-agent"

agent_dir="/usr/local/webcron"

remote_cmd="ssh $USER@$SERVER"

echo "start to deploy webcron agent"

go build -o $FILE

$remote_cmd "if ! test -d $agent_dir; then mkdir -p $agent_dir; fi"
scp $FILE $USER@$SERVER:$agent_dir

echo "upload to serverv $SERVER"
$remote_cmd "cd $agent_dir;source ./$FILE &"

echo "webcron agent service started success"



