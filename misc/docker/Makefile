SHELL:=/bin/bash

build: git
	sudo docker build -t orca -f Dockerfile ../..

run: build
	sudo docker run --env "ORCA_MYSQL=orca:orca@tcp(172.17.0.1:3306)/orca" \
		--mount type=bind,source=/tmp/orca/tmp,target=/tmp \
		--mount type=bind,source=/tmp/orca/server/dist/prod/log,target=/log \
		orca

install: build
	sudo docker save orca > docker-orca.tar
	scp docker-orca.tar freya:/tmp
	scp install.sh freya:/tmp
	ssh freya "chmod +x /tmp/install.sh && /tmp/install.sh && rm /tmp/install.sh"

git:
	../../server/build/git-hash.sh > ../../server/build/.git-hash

clean:
	sudo docker stop orca || :
	sudo docker rm orca || :
	sudo docker rmi orca || :
	sudo docker image prune --force
