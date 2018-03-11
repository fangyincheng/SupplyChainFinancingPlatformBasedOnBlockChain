#!/bin/bash

set -e

#设置路径到脚本目录
basepath=$(cd `dirname $0`; pwd)
cd $basepath

#设置docker-compose.yaml使用的环境变量
if [ -z $1 ]; then
    export IMAGE_TAG=latest
else
    export IMAGE_TAG=$1
fi
export COMPOSE_PROJECT_NAME=supply_chain

#卸载网络
docker-compose -f docker-compose.yaml down