#!/bin/bash
while true; do

inotifywait -e modify,create,delete -r ./ && \
	clear
	go fmt ./... \
		&& go build -o reverse-proxy \
		&& go test  -count=1 ./... 
done
