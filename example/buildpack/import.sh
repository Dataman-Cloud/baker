#!/bin/bash

../../bin/baker -s 192.168.1.21:8000 buildpack import --name app --from tomcat:8.0 --binaryFile sample.zip --binaryPath /usr/local/tomcat/webapps --startCmd catalina.sh\ run --disconf true  
