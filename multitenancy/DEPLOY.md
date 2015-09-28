#Deploy locally
```
# Install and run docker
  1. brew install boot2docker
  2. boot2docker init
  3. boot2docker up
  4. edit server/z-manager.sh
    a) DOCKER_HOST="tcp://<your docker ip>:<your docker port>"
    b) comment unset DOCKER_TLS_VERIFY
    c) comment unset DOCKER_CERT_PATH
  5. $(boot2docker shellinit)`
  6. docker pull nflabs/zeppelin-bhs-aug-5-spark-1.4-hadoop-2.0.0-cdh-4.7.0:latest

# Install and run spark
  1. brew install apache-spark
  2. /usr/local/Cellar/apache-spark/1.4.1/libexec/sbin/start-master.sh
  3. open localhost:8080 in the browser and copy spark hostname
  4. /usr/local/Cellar/apache-spark/1.4.1/libexec/sbin/start-slave.sh spark://<your hostname>:7077
  5. edit server/z-manager.sh
    a) SPARK_URL="http://localhost:8080”

# Start manager
  1. ./server/z-manager.sh
```
