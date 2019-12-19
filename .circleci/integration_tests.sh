#!/bin/sh

set -e pipefail


#############
## globals ##
#############


cleanup() {
	echo "kill processes"
	pkill reverse-proxy
	pkill passthrough-service
}
trap cleanup EXIT

#########
## env ##
#########

export PORT="9001"
export services="passthrough,example"
export passthrough_incoming_path="/passthrough/"
export passthrough_outgoing_url="http://localhost:9002/"
export example_incoming_path="/example/"
export example_outgoing_url="http://example.com"

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
URL="http://localhost:9001/passthrough/ping?pause=1"
echo "making request to: $URL"
wget -O- $URL
echo "proxy log: "
cat proxy.log

> proxy.log
URL="http://localhost:9001/example/"
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