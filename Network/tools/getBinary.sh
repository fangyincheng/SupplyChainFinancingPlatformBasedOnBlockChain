#!/bin/bash

set -e

echo "===> 设置路径到脚本目录"
basepath=$(cd `dirname $0`; pwd)
cd $basepath

echo "===> 设置版本和系统参数环境变量，用于下一步的下载链接"
if [ -z $1 ]; then
    export VERSION=1.0.0
else
    export VERSION=$1
fi
export ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')

echo "===> 下载二进制工具(VERSION=${VERSION}, ARCH=${ARCH})"
curl https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/${ARCH}-${VERSION}/hyperledger-fabric-${ARCH}-${VERSION}.tar.gz | tar xz
