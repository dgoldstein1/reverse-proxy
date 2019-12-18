#!/bin/sh

set -e pipefail

#########
## env ##
#########

export PORT="9001"
export services="google,wikipedia"

export google_incoming_path="/google/"
export google_outgoing_url="http://google.com"

export wikipedia_incoming_path="/wikipedia/"
export wikipedia_outgoing_url="http://wikipedia.org"

################
## run binary ##
################

ls reverse-proxy
./reverse-proxy > proxy.log 2>&1 &
pid=$!
> proxy.log
sleep 1
cat proxy.log

###############
## run tests ##
###############


> proxy.log
URL="http://localhost:$PORT$google_incoming_path"
echo "making request to: $URL"
curl $URL
cat proxy.log


> proxy.log
URL="http://localhost:$PORT$wikipedia_incoming_path"
echo "making request to: $URL"
curl $URL
cat proxy.log

##############
## clean up ##
##############

echo "stopping reverse-proxy on pid $pid"
kill $pid