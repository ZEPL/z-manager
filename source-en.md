## To run Z-Manager

Paste this into a Terminal prompt.
```shell
curl -fsSL https://raw.githubusercontent.com/NFLabs/z-manager/master/zeppelin-installer.sh | bash
```

The script explains what it will do and then pauses before doing it.


##
### What is Z-Manager?
A simple tool that automates process of getting Zeppelin up and running

**Value proposition:**
  - *No Build*

    No need to build Zeppelin manually, we distribute up-to-date convenience binaries for you.

  - *No Configuration*

    No need for manual configuration. Just pick Hadoop and Spark version though the installation.

### Who should use it?
Anybody who want to try Zeppelin and save some time on building\configuration process should

### Why?
> Installing and configuring Zeppelin is a laborious process: you need to build
  Zeppelin using [right set of parameters]() that depends on your setup.
  Z-Manager allows you to skip it by asking few questions and downloading everything.

### When?
> Either for trying Zeppelin with Spark  (no need for separate Spark installation or existing cluster) or before plugging Zeppelin to your cluster


Please mind that it's in early stage right now.




### UI: interactive CLI
Z-Manager have a simple interactive command line interface:
<div id="video"></div>

### UI: non-interactive CLI
In case you want to install Zeppelin in a non-interactive fashion, run `./zeppelin-installer.sh -h` and provide the options though arguments.


## What is happening underneath?

### Download latest binary Zeppelin build

Z-Manager download and unpacks [Apache Zeppelin (incubating)](zeppelin.incubator.apache.org) on your machine.
By picking right spark and hadoop version and downloading binary it frees you from building Zeppelin from sources.


### Configure Zeppelin

Configure Spark cluster mode using one of the following presets
  - Standalone
  - YARN
  - Mesos



### See what's installed

Z-Manager will install latest Zeppelin in current directory and touch nothing else.
Then you can move a Zeppelin installation wherever you like.

```shell
$ls .
  zeppelin-0.5.0-incubating-SNAPSHOT
  
$cd zeppelin-0.5.0-incubating-SNAPSHOT
$find .
  ./bin/zeppelin-daemon.sh
  ./conf/zeppelin-env.sh.template
  ./conf/zeppelin-site.xml
  ./conf/zeppelin-site.xml.template
  ./interpreter/angular/zeppelin-angular-0.5.0-incubating-SNAPSHOT.jar
  ./interpreter/hive/zeppelin-hive-0.5.0-incubating-SNAPSHOT.jar
  ./interpreter/md/zeppelin-markdown-0.5.0-incubating-SNAPSHOT.jar
  ./interpreter/sh/zeppelin-shell-0.5.0-incubating-SNAPSHOT.jar
  ./interpreter/spark/zeppelin-spark-0.5.0-incubating-SNAPSHOT.jar
  ...
```



## Further work
We are planning a full featured GUI version of manager soon, with more external integrations and pre-sets

