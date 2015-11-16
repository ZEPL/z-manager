# Z-Manager
> Simplify getting Zeppelin up and running

## What is Z-Manager
Z-Manager is a set of tools that simplify deployment of Apahce Zeppelin.

It consists of 2 main parts:
  - **Installer** a single click installation and configuration tool of Zeppelin with Spark integration.

  It supports external clusters with Standalone, Mesos and YARN resource managers.

  - **Multi-tenancy** *(New!)*

  Allows you to run a containerized instance of Zeppelin, one per-user with exteranl authentification

In case of only issues or question please post on the  mailing list [z-manager@googlegroups.com](https://groups.google.com/forum/#!forum/z-manager) or create an Issue here.

## Multi-tenancy
A single binary app with reverse-proxy to the container fleet of Zeppelin instances.
See the `./multitenancy/README.md`

## Installer

**Interactive CLI UI**
```
  curl -fsSL https://raw.githubusercontent.com/NFLabs/z-manager/master/zeppelin-installer.sh | bash
```

**Non-interatcive CLI**
```
./zeppelin-installer.sh -h
```

## Hapoop\Spark compatibility

|                      |  1.2.1 |  1.3.0 | 1.3.1|1.4.0  |1.4.1 |
| -------------------- | :----: | :----: | :---:|:-----:|:----:|
|  1.x                 |        |        |      |       |   x  |
|  2.3                 |        |        |      |       |      |
|  2.4 and later       |        |        |   x  |       |      |
|  2.6                 |        |        |      |   x   |      |
|  2.7.1               |        |        |      |       |   x  |
|  CDH4: 2.0.0-cdh4.7.1|        |    x   |   x  |   x   |
|  CDH5: 2.5.0-cdh5.3.0|        |        |   x  |       |
|  HDP?                |        |        |      |       |
|  MapR 3.x            |        |        |      |       |
|  MapR 4.x            |        |        |      |       |



## Example
![Video of Z-Manager](https://raw.githubusercontent.com/NFLabs/z-manager/master/yarn.gif)


**Disclaimer**

>Z-Manager does not collect any personal information.
>In order to measure the efforts and assure the best product quality, NFLabs reserves the right to count the number of Zeppelin installations made through Z-Manager.
>You can always avoid that by using `--no-count` CLI option.
