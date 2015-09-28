#!/bin/bash
#
# Script to run Z-Manager with configuration
#
# Bluehole cluster configuration

export SERVE_WEBAPP_FROM_FS='false'
export MARTINI_ENV='production'
#Use this to change port
export PORT="80"

export SPARK_URL="http://bhsmaster2.nflabs.com:8080"
export HUB_URL="http://bhsmaster1.nflabs.com:8042"

export USERS_FOLDER_PATH="/data/users"
export DEFAULT_USER_FOLDER_NAME="defaultuser"

export DOCKER_IMAGE="nflabs/zeppelin-bhs-sep-25-with-zeppelinhub-integration-spark-1.4-hadoop2.0.0-mr1-cdh-4.7.1:latest"

#All BHS cluster uses port :2375, see _configs/docker
export DOCKER_PORT="2375"

export DOCKER_HOSTS="bhsworker1: mina,bird
bhsworker2: fbdkdud93,ipsae,bhpa
bhsworker3: khalidhuseynov,suomi,aronsate
bhsworker4: anthonycorbacho,bht2,kiriejs
bhsworker5: heidi,thenewknow
bhsworker6: astroshim,cris,sung
bhsworker7: hifive7,mssk26
bhsworker8: admin,bhh5
bhsworker9: bhfs,mdi_admin
bhsworker10:"


echo "Starting a Z-Manager server"

./server "$@" 2>&1 > z-manager.log
