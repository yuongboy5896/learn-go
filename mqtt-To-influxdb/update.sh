#!/bin/bash
URL=$1
go build -o mqtt2influxdb .
export IMAGEURL=${URL}
TIME=`date "+%Y%m%d%H%M"`

IMAGE_NAME=${IMAGEURL}/common/mqtt2influxdb:${TIME}
DEP="mqtt2influxdb"
NAME="mqtt2influxdb"
echo $IMAGE_NAME

docker build -t ${IMAGE_NAME} .
docker push  ${IMAGE_NAME}
docker rmi ${IMAGE_NAME}