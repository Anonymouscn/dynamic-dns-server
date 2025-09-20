#!/bin/bash

bash clean.sh
# 交叉编译构建镜像
docker buildx build --platform linux/amd64,linux/arm64 -t pgl888999/dynamic-dns-server --push .
# 打包 secret 文件夹
PASSWD=$RANDOM
echo "secret 随机密码: ${PASSWD}"
zip -rP ${PASSWD} ./bin/secret.zip ./secret