# builder stage

FROM golang:latest as builder

ARG DOCKER_RUNNING=yes
ARG DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -yq protobuf-compiler 
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0

COPY . /project

RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	/project/server/build/build-server.sh prod

# clean stage

FROM alpine:latest

RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /project/server/dist/prod/orca-server-next /orca-server

COPY ./misc/docker/http-ca.crt /etc/ssl/certs/elastic-http-ca.crt

ENV TZ="Asia/Shanghai"

ENV ORCA_WEB=":80"
ENV ORCA_LOG="/log"
ENV ORCA_DIR="/tmp"

CMD ["/orca-server"]
