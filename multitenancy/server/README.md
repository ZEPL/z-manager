# Z-Manager


## Build

### Frontend webapp

See `../web/ui/README.md` for details.

*tl;dr*

```
cd ../web/ui
npm install
npm run build
cd -
```

You have to **build for production once**, for the first time

### Produciton build
Serving resources from the binary itself

```
export GOPATH=$PWD/third_party/go:$GOPATH
#only once, to get go-bindata tool
go get -u github.com/jteeuwen/go-bindata/...

cd server
# to generat binay assets
$GOPATH/bin/go-bindata -o ./assets.go ../web/ui/public/**
go build ./...
```


### Dev build
Serving resources from the filesystem, under `../web/ui/public`

```
export GOPATH=$PWD/third_party/go:$GOPATH
cd server
go build ./...
```


### Autoupdate server
Serves binary diff updates from https://github.com/bzz/test/releases
```
#only once, to get autoupdate-server binary
go get -u github.com/bzz/autoupdate-server/...
autoupdate-server -k private.pem -o bzz -n zeppelin-manager
#or from clone
#go run *.go -k private.pem  -o bzz -n autoupdate-server
```

## Cross-complie Z-Manager for Linux
```
GOOS=linux GOARCH=386 go build ./...
```


## Run

To run z-manager server use the provided convenience .sh wrapper `./z-manager.sh`

Check it for example of configuration parameters that need to be set for server to be operational: url to Hub, Spark Master, Docker instances, etc

**important** to run as a load-balancer for the cluster of docker instances, *each Docker machine* must:
   - have a `../_config/find_open_port.py` service running
   - have a docker deamon, listening on TCP port `DOCKER_PORT` (i.e in `../_config/docker` we use `-H "tcp://$(hostname --ip-address):2375"`)

For the first time "Create" button will fetch docker image of Zeppelin `DOCKER_IMAGE`, it will take for a while to download ~1Gb, so please be patient.


## Build&Publish Docker image of Zeppelin

    # build zeppelin
    mvn clean package -Phadoop-2.2 -Dhadoop.version=2.0.0-cdh4.7.1 -Pspark-1.4 -Dspark.version=1.4.0 -Pyarn -Ppyspark -DskipTests -Pbuild-distr
    # copy final distr to .
    cp /<path to your zeppelin>/zeppelin-distribution/target/zeppelin-0.6.0-incubating-SNAPSHOT.tar.gz .
    # build the image
    docker build -t nflabs/zeppelin-bhs-aug-5-spark-1.4-hadoop-2.0.0-cdh-4.7.0 .

    # publsh it on DockerHub
    ## login to DockerHub
    docker login

    ## push
    docker push nflabs/zeppelin-bhs-aug-5-spark-1.4-hadoop-2.0.0-cdh-4.7.0

    ## on all workers
    docker pull nflabs/zeppelin-bhs-aug-5-spark-1.4-hadoop-2.0.0-cdh-4.7.0


## Add a new user

Z-Manager uses ZeppelinHub for authentification, so you need to
  - create a user on a Hub instance from `HUB_URL` (set in `./z-manager.sh`)
  - add his name to `DOCKER_HOSTS` in `./z-manager.sh`, assiging user to the particular host


## Bluehole AKA BHS deployment

Z-Manager is running on a `bhsmaster3.nflabs.com:80` from `/home/root`, using a `./z-manager-bhs.sh` script.
Please feel free to move it but update this document accordingly.

It had WebUI + reverse proxy, so every user's containers is created on remote machine with Docker daemon and all user requests are forwarded there. All docker machines should be using modified daemon script from `_config/docker` as it enables remote interaction (off by default)

All the user's notebooks and zeppelin instance configuration is saved on the shaared NFS under `/data/users/<user name>` with `/data/users/defaultuser` being a template, used for a new users (on the first login)

## Localmode deployment
  - Docker REST API daemon should be enabled on "tcp://$DOCKER_HOST:$DOCKER_PORT"
  - Spark should be up and running, and Spark url is passed through `SPARK_URL` environment variable. It's possible to run Spark locally (e.g. http://localhost:8080) as well as remotely.
  - `USERS_FOLDER_PATH` points to `_/config/users/`
