# builder stage

FROM golang:latest as builder

ARG DOCKER_RUNNING=yes
ARG DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -yq protobuf-compiler tzdata ca-certificates
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0

COPY . /project

RUN /project/server/build/build-server.sh prod

# clean stage

FROM alpine

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/
COPY --from=builder /project/server/dist/prod/orca-server-next /orca-server

COPY ./misc/docker/http-ca.crt /etc/ssl/certs/elastic-http-ca.crt

ENV TZ="Asia/Shanghai"

ENV ORCA_WEB=":80"
ENV ORCA_LOG="/log"
ENV STATIC_DIR="/tmp"

CMD ["/orca-server"]
