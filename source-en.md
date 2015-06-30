## To run Z-Manager

Paste this into a Terminal prompt.
```shell
curl -fsSL https://raw.githubusercontent.com/NFLabs/z-manager/master/zeppelin-installer.sh | bash
```

The script explains what it will do and then pauses before doing it.


##
<a href="#what-is-z-manager-">
### What is Z-Manager?
</a>
A simple tool that automates process of getting Zeppelin up and running

<br/>
**Value proposition:**
  - *No Build*

    No need to build Zeppelin manually, we distribute up-to-date convenience binaries for you.

  - *No Configuration*

    No need for manual configuration. Just pick Hadoop and Spark version though the installation.


<a href="#who-should-use-it-">
### Who should use it?
</a>
Anybody who want to try Zeppelin and save some time on building\configuration process should

### Why?
> Installing and configuring Zeppelin is a laborious process: you need to build
  Zeppelin using [right set of parameters](http://zeppelin.incubator.apache.org/docs/install/install.html) that depends on your setup.
  Z-Manager allows you to skip it by asking few questions and downloading everything.

### When?
> Either for trying Zeppelin with Spark  (no need for separate Spark installation or existing cluster) or before plugging Zeppelin to your cluster


Please mind that it's in early stage right now.



<a href="#ui-interactive-cli">
### UI: interactive CLI
</a>
Z-Manager have a simple interactive command line interface:
<div id="video"></div>

<a href="#ui-non-interactive-cli">
### UI: non-interactive CLI
</a>
In case you want to install Zeppelin in a non-interactive fashion, run `./zeppelin-installer.sh -h` and provide the options though arguments.


## What is happening underneath?

### Download latest binary Zeppelin build

Z-Manager download and unpacks [Apache Zeppelin (incubating)](zeppelin.incubator.apache.org) on your machine.
By picking right spark and hadoop version and downloading binary it frees you from building Zeppelin from sources.

<a href="#configure-zeppelin">
### Configure Zeppelin
</a>

Configure Spark cluster mode using one of the following presets
  - Standalone
  - YARN
  - Mesos


<a href="#see-what-s-installed">
### See what's installed
</a>

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

<a href="#what-is-z-manager-">
### Feedback & Questions
</a>
If you have any comments, questions, or suggestions for Z-Manager please feel free to file a [new issue](https://github.com/NFLabs/z-manager/issues).

## Further work
We are planning a full featured GUI version of manager soon, with more external integrations and pre-sets

