#!/bin/bash

#Remove containers
docker rm -f $(docker ps -aq)

#Prune volumes
docker volume prune

#Prune network
docker network prune

docker rmi -f $(docker images | awk '($1 ~ /dev-peer.*.walletdemo.*/) {print $3}')

#Remove certificates
#rm -rf crypto-config/

#Remove channel artifacts
#rm -rf channel-artifacts/*
