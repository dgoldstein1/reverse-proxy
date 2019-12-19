#!/bin/sh

set -e pipefail


#############
## globals ##
#############

cleanup() {
	pkill reverse-proxy
	pkill passthrough-service
}
trap cleanup EXIT

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
> proxy.log
sleep 1
echo "proxy logs:"
cat proxy.log


###############################
## start passthrough service ## 
###############################


ls .circleci/passthrough-service
export PORT=9002
.circleci/passthrough-service > passthrough.log 2>&1 &
sleep 1
wget -O- -q localhost:9002/ping
echo "passthrough log:"
cat passthrough.log


###############
## run tests ##
###############

> proxy.log
URL="http://localhost:9001/passthrough/ping"
echo "making request to: $URL"
wget -O- $URL
echo "proxy log: "
cat proxy.log

##############
## success! ##
##############

echo "============="
echo "== SUCCESS =="
echo "============="