#!/bin/bash

../../bin/baker buildpack import --name app --from tomcat:8.0 --binaryFile sample.zip --binaryPath /usr/local/tomcat/webapps
