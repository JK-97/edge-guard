#!/bin/bash

### BEGIN INIT INFO
# Provides:          jiangxing team
# Required-Start:    $local_fs $network
# Required-Stop:     $local_fs
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: jiangxing team
# Description:       jiangxing core daemon
### END INIT INFO
while true
do
cd /edge/jxcore/bin
/edge/jxcore/bin/jxcore serve 1>> /edge/logs/jxcore-stdout-`date +"%Y-%m-%d"`.log 2>> /edge/logs/jxcore-stderr-`date +"%Y-%m-%d"`.log 
done
#exit 0
