#!/bin/bash
while true; do

inotifywait -e modify,delete -r ./ && \
	clear
	go fmt ./... \
		&& go build -o reverse-proxy \
		&& go test  -count=1  -coverprofile=coverage.out ./... \
		&& .circleci/integration_tests.sh 
done
 \