#!/bin/bash
#
# Script to run Z-Manager
#
# It is optimized for local dev mode on osx: same machine running
# docker (boot2docker), spark and z-manager.
# It severs weapp from ../webapp/ui/public, and expects $(boot2docker shellinit)

export SERVE_WEBAPP_FROM_FS='true'

#Use this to change port
#export PORT="3001"

export SPARK_URL="http://localhost:8080"
export HUB_URL="http://dev.zeppelinhub.com"

#absolute path on Host FS
export USERS_FOLDER_PATH="$(cd "$(dirname "../_configs/users/.")"; pwd)/$(basename "$1")"
export DEFAULT_USER_FOLDER_NAME="defaultuser"

#use this if on linux, without TLS
#export DOCKER_HOST='localhost'
#unset DOCKER_TLS_VERIFY
#unset DOCKER_CERT_PATH


export DOCKER_IMAGE="nflabs/zeppelin-bhs-spark-1.4-hadoop-2.0.0-mr1-cdh-4.2.0:latest"
export DOCKER_PORT="2376"

export DOCKER_HOST_1="${DOCKER_HOST##*/}" #get only the hostname
export DOCKER_HOSTS="${DOCKER_HOST_1%%:*}: khalidhuseynov, anthonycorbacho"


if [[ ! -f "./server" ]]; then
  echo "No built server found, building it.."

  if [[ ! -d "../web/ui/public" ]]; then
    echo "  No built client web-app found, building it.."
    pushd '../web/ui'
    npm install
    npm run build
    popd
    echo "  Done"
  fi

  go-bindata -o ./assets.go ../web/ui/public/**
  go build ./...
  echo "Done"
fi

echo "Starting a Z-Manager server"
./server "$@" 2>&1 | tee -a z-manager.log
