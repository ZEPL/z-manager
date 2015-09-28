FROM centos:centos6
MAINTAINER NFLabs <contacts@nflabs.com>
# requites built Zeppelin distribution (-Pbuild-distr) in the current directory
  
# Update the image with the latest packages
RUN yum update -y; yum clean all

# Get wget
RUN yum install wget -y; yum clean all

# Remove old jdk
RUN yum remove java; yum remove jdk

# Install oracle jdk7
RUN wget --continue --no-check-certificate --header "Cookie: oraclelicense=accept-securebackup-cookie" -O jdk-linux-x64.rpm "http://download.oracle.com/otn-pub/java/jdk/7u51-b13/jdk-7u51-linux-x64.rpm"
RUN rpm -Uvh jdk-linux-x64.rpm
RUN rm jdk-linux-x64.rpm

ENV JAVA_HOME /usr/java/default
ENV PATH $PATH:$JAVA_HOME/bin

# Set zeppelin env
ENV ZEPPELIN_NOTEBOOK_DIR /zeppelin/notebook

RUN mkdir /usr/lib/zeppelin
ADD zeppelin-0.6.0-incubating-SNAPSHOT.tar.gz /tmp/
RUN cp -rf /tmp/zeppelin-0.6.0-incubating-SNAPSHOT/* /usr/lib/zeppelin

# Get mysql client to access remote Hive Metastore
RUN wget -O /usr/lib/zeppelin/interpreter/spark/mysql-connector-java.jar http://search.maven.org/remotecontent?filepath=mysql/mysql-connector-java/5.1.26/mysql-connector-java-5.1.26.jar 


# Change timezone to Seoul
RUN ln -sf /usr/share/zoneinfo/Asia/Seoul /etc/localtime

# Open docker port 8080
EXPOSE 8080
