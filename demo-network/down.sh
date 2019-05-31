#!/bin/bash

#Remove containers
docker rm -f $(docker ps -aq)

#Prune volumes
docker volume prune

#Prune network
docker network prune

#Remove certificates
#rm -rf crypto-config/

#Remove channel artifacts
rm -rf channel-artifacts/*
