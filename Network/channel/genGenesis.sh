#!/bin/bash

set -e

echo "===> 设置路径到脚本目录"
basepath=$(cd `dirname $0`; pwd)

echo "===> 设置二进制工具路径和执行工具所需配置路径的环境变量"
export PATH=$basepath/../tools/bin:$PATH
export FABRIC_CFG_PATH=$basepath/testchainid

echo "===> 执行工具生成genesis.block"
configtxgen -profile OrdererGenesis -outputBlock $basepath/testchainid/genesis.block