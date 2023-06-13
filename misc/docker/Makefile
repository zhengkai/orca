SHELL:=/bin/bash

include ../../server/build/env.sh

build: git
	sudo docker build -t orca -f Dockerfile ../..

run: build
	sudo docker run \
		--env "OPENAI_API_KEY=$(OPENAI_API_KEY)" \
		--env "ORCA_ES_ADDR=$(ORCA_ES_ADDR)" \
		--env "ORCA_ES_USER=$(ORCA_ES_USER)" \
		--env "ORCA_ES_PASS=$(ORCA_ES_PASS)" \
		--mount type=bind,source=/www/orca/static,target=/tmp \
		--mount type=bind,source=/www/orca/server/dist/prod/log,target=/log \
		-p 127.0.0.1:21035:80 \
		orca

install: build
	sudo docker save orca > docker-orca.tar
	scp docker-orca.tar lamia:/tmp
	scp install.sh lamia:/tmp
	scp ../../server/build/env.sh lamia:/tmp
	ssh lamia "chmod +x /tmp/install.sh && /tmp/install.sh && rm /tmp/install.sh"
	ssh lamia "rm /tmp/env.sh"

nginx:
	scp ../nginx/prod.conf lamia:/etc/nginx/vhost.d/600-orca

git:
	../../server/build/git-hash.sh > ../../server/build/.git-hash

clean:
	sudo docker stop orca || :
	sudo docker rm orca || :
	sudo docker rmi orca || :
	sudo docker image prune --force
