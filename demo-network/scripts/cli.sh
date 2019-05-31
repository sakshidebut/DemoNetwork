#!/bin/bash

#Create and join channel
export CHANNEL_NAME=demochannel
peer channel create -o orderer.india.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/india.com/orderers/orderer.india.com/msp/tlscacerts/tlsca.india.com-cert.pem

#Join peer to channel for peer0 org1
peer channel join -b demochannel.block
#Install chaincode
peer chaincode install -n democc -v 1.0 -p github.com/chaincode/mychaincode/go/

#Instantiate chaincode
peer chaincode instantiate -o orderer.india.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/india.com/orderers/orderer.india.com/msp/tlscacerts/tlsca.india.com-cert.pem -C $CHANNEL_NAME -n democc -v 1.0 -c '{"Args":["init"]}' -P "AND ('Org1MSP.peer')"

#Add New User
#peer chaincode invoke -o orderer.india.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/india.com/orderers/orderer.india.com/msp/tlscacerts/tlsca.india.com-cert.pem  -C $CHANNEL_NAME -n democc --peerAddresses peer0.org1.india.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.india.com/peers/peer0.org1.india.com/tls/ca.crt -c '{"Args":["createUser","Sakshi","sakshi.bansal@debutinfotech.com","9876543210"]}'

#Query Specific User
#peer chaincode invoke -o orderer.india.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/india.com/orderers/orderer.india.com/msp/tlscacerts/tlsca.india.com-cert.pem  -C $CHANNEL_NAME -n democc --peerAddresses peer0.org1.india.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.india.com/peers/peer0.org1.india.com/tls/ca.crt -c '{"Args":["getWalletBalance","2056000e4e520373af072c5ca58a65c7b758b5f063a03d7c82596d9b3458544f"]}'
