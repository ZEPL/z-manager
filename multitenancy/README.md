# Zeppelin Multitenancy

> Makes using Zeppelin w/ ZeppelinHub easy in multi-user environemnt

## Problem
In organization it is not feaseble to make all users install and configure Zeppelin.

This application allows you to have a single URL where all users go to get their own, containerized version of Zeppelin, utilasing resources of the shared cluster.

## Solution
![Multitenancy architecture](https://raw.githubusercontent.com/NFLabs/z-manager/master/multitenancy/architecture.png)

*This is Beta, for now only Spark Standalone is supported*
  

## Build
Build a React frontend webapp, see `web/ui/READEME.md` for details

Build a Golang backend webserver, see `server/README.md`

**tl;dr**
`./build_all.sh`

## Setup it on your cluster

All configureation for now is kept in the `./server/z-manager.sh`.
See the `./server/z-manager-on-your-cluster.sh` as an example of cluster configuration.

**Pre-requests**

  * cluster of Docker machines, with HTTP API enabled (see example of docker deamon script, updated for HTTP in  `./_config/dcoker`)
  * each machine have `USERS_FOLDER_PATH` mounted to the same location (shared Notebook storage, each user will have it's own dir)
  * each machine have `./_config/find_open_port.py` running (free port discovery)
  * `./z-manager.sh` configured w/ user accounts, pointing to Docker and a SparkMaster URL


## Details

Essentially, Z-Manager multitenancy is very simple load-balancer and reverse-proxy with authentification for:
 - HTTP trafic to zeppelin instances (one per user), run inside Docker containers on one or many machines
 - Websocket traffic to same zeppelin instances
 - HTTP traffic to SparkMaster, to get cluster status\available resources
