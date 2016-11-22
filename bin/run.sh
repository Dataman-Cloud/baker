#! /bin/bash

#------------------------------------------------------
# CHECK ENVIRONMENT VALUE: CONFIG_DIR AND CONFIG_SERVER
#------------------------------------------------------
if [ ! $CONFIG_SERVER ] || [ ! $CONFIG_DIR ]; then 
	echo "no setting for config server and config dir. Please set parameter in docker run."
	exist 1
fi

#------------------------------------------------
# DOWNLOAD CONFIG FILES
#------------------------------------------------
mkdir -p /config
./baker -s $CONFIG_SERVER disconf pull --path=$CONFIG_DIR 
./baker -s $CONFIG_SERVER disconf unzip --file=props.zip --path=/config
mv /config/*/ /
