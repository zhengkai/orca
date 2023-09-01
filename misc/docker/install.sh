#!/bin/bash

TARGET="Lamia"

if [ "$HOSTNAME" != "$TARGET" ] && [ "$HOSTNAME" != "Rachel" ]; then
	>&2 echo only run in server "$TARGET"
	exit 1
fi

cd "$(dirname "$(readlink -f "$0")")" || exit 1
if [ ! -e ./env.sh ]; then
	>&2 echo no env file
	exit 1
fi

. ./env.sh || exit 1

sudo docker stop orca
sudo docker rm orca
sudo docker rmi orca

sudo cat /tmp/docker-orca.tar | sudo docker load

sudo docker run -d --name orca \
	--env "OPENAI_API_KEY=${OPENAI_API_KEY}" \
	--env "ORCA_ES_ADDR=${ORCA_ES_ADDR}" \
	--env "ORCA_ES_USER=${ORCA_ES_USER}" \
	--env "ORCA_ES_PASS=${ORCA_ES_PASS}" \
	--env "ORCA_VA_TOKEN=${ORCA_VA_TOKEN}" \
	--mount type=bind,source=/www/orca/log,target=/log \
	--mount type=bind,source=/www/orca/static,target=/tmp \
	-p 127.0.0.1:21035:80 \
	--restart always \
	orca
