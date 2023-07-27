#!/usr/bin/env bash

VAR_FR_DOCKER_IMAGE_NAME="financial_statement"
VAR_FR_DOCKER_IMAGE_VERSION=""

function PushImage() {
    eval "docker login registry.intsig.net"
    eval "docker tag ${VAR_FR_DOCKER_IMAGE_NAME}:${VAR_FR_DOCKER_IMAGE_VERSION} registry.intsig.net/textin_com/${VAR_FR_DOCKER_IMAGE_NAME}:${VAR_FR_DOCKER_IMAGE_VERSION}"
    eval "docker push registry.intsig.net/textin_com/${VAR_FR_DOCKER_IMAGE_NAME}:${VAR_FR_DOCKER_IMAGE_VERSION}"
    if [[ $? -eq 0 ]];
    then
        echo 'finished!'
    else
        echo 'push failed!'
    fi
}

function DockerBuild () {
    read -e -p "输入docker image 的版本号：" -i "2.0.0" VAR_FR_DOCKER_IMAGE_VERSION
    eval "docker build -t ${VAR_FR_DOCKER_IMAGE_NAME}:${VAR_FR_DOCKER_IMAGE_VERSION} -f Dockerfile ."
    if [[ $? -eq 0 ]];
    then
        return 0
    else
        return 1
    fi
}


echo '菜单'
echo '1.编译Docker Image'
read -p '请选择：' cmd

if [[ $cmd == "1" ]]; then
    DockerBuild
  if [[ $? -eq 0 ]]; then
    read -e -p "构建成功，是否要推送Image：${VAR_FR_DOCKER_IMAGE_NAME}:${VAR_FR_DOCKER_IMAGE_VERSION} 到 registry.intsig.net/textin_com/${VAR_FR_DOCKER_IMAGE_NAME}:${VAR_FR_DOCKER_IMAGE_VERSION} ：" -i "y" needPushImage
    if [[ $needPushImage == "y" ]]; then
        PushImage
    fi
  else
    echo '构建失败！'
  fi
fi






