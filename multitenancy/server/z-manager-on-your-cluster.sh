#!/bin/bash
#
# Script to run Z-Manager with configuration
#
# Your own cluster configuration

export SERVE_WEBAPP_FROM_FS='false'
export MARTINI_ENV='production'
#Use this to change port
export PORT="80"

export SPARK_URL="http://sparkmaster.yourcompany.com:8080"
export HUB_URL="http://zeppelinhub.yourcompany.com:8042"

#shared FS, should be availabele on docker cluster for storing user notebooks/configuration
export USERS_FOLDER_PATH="/data/users"
export DEFAULT_USER_FOLDER_NAME="defaultuser"

export DOCKER_IMAGE="<your-zeppelin-docker-image>:latest"

#Port for docker HTTP :2375, see _configs/docker
export DOCKER_PORT="2375"

#point to the docker cluster
export DOCKER_HOSTS="server1.yourcompany.com: username1, username2
server2.yourcompany.com: username3, username4
server3.yourcompany.com: username5, username6
server4.yourcompany.com: username7, username8
server4.yourcompany.com:
"


echo "Starting a Z-Manager server"

./server "$@" 2>&1 > z-manager.log
