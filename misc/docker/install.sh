#!/bin/bash

TARGET="Freya"

if [ "$HOSTNAME" != "$TARGET" ]; then
	>&2 echo only run in server "$TARGET"
	exit 1
fi

sudo docker stop orca
sudo docker rm orca
sudo docker rmi orca

sudo cat /tmp/docker-orca.tar | sudo docker load

sudo docker run -d --name orca \
	--env "TANK_MYSQL=orca:orca@tcp(172.17.0.1:3306)/orca" \
	--env "STATIC_DIR=/tmp" \
	--env "OUTPUT_PATH=/output" \
	--mount type=bind,source=/www/orca/output,target=/output \
	--mount type=bind,source=/www/orca/log,target=/log \
	--mount type=bind,source=/www/orca/static,target=/tmp \
	--restart always \
	orca
