#!/bin/sh

set -e pipefail

#########
## env ##
#########

export PORT="9001"
export services="passthrough"

export passthrough_incoming_path="/passthrough/"
export passthrough_outgoing_url="http://localhost:9002/"

################
## run binary ##
################

ls reverse-proxy
./reverse-proxy > proxy.log 2>&1 &
pid=$!
> proxy.log
sleep 1
cat proxy.log


###############################
## start passthrough service ##
###############################


ls .circleci/passthrough-service
export PORT=9002
.circleci/passthrough-service > passthrough.log 2>&1 &
pid_passthrough=$!


###############
## run tests ##
###############

> proxy.log
URL="http://localhost:9001/passthrough/ping"
echo "making request to: $URL"
wget -O- $URL
cat proxy.log

##############
## clean up ##
##############

echo "stopping reverse-proxy on pid $pid"
kill $pid

echo "stopping passthrough service on pid $pid_passthrough"
kill $pid_passthrough


echo "============="
echo "== SUCCESS =="
echo "============="