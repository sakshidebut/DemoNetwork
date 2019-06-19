#!bin.bash
export CHANNEL_NAME=mychannel

peer chaincode install -n walletdemo -v $1  -p github.com/chaincode/demo-network/cmd/

peer chaincode upgrade -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n walletdemo -v $1 -c '{"Args":["initLedger"]}' -P "OR ('Org1MSP.peer')"