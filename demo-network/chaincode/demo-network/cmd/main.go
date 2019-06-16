// Package main Implements the Init & Invoke functions, Starts the chaincode
package main

import (
	"fmt"

	"github.com/chaincode/demo-network/pkg/users"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
)

// Chaincode default chaincode implementation with router
type Chaincode struct {
	router *router.Group
}

// Init initializes chain code - sets chaincode "owner"
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.HandleInit(stub)
}

// Invoke - entry point for chain code invocations
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}

// New Define the Router
func New() *Chaincode {
	// create a new router instance
	r := router.New("Chaincode")
	chaincode := &Chaincode{r}

	// Handle the init/upgrade
	r.Init(invokeInit)

	// Other routes

	/***** users routes *****/

	r.Invoke(`getUser`, users.GetUser, param.Struct(`data`, &users.User{}))
	r.Invoke(`getUsers`, users.GetUsers, param.Struct(`data`, &users.UserId{}))
	r.Invoke(`getAssets`, users.GetAssets, param.Struct(`data`, &users.UserId{}))
	r.Invoke(`addAsset`, users.AddAsset, param.Struct(`data`, &users.Asset{}))
	r.Invoke(`checkAsset`, users.CheckAsset, param.Struct(`data`, &users.CheckAssetStruct{}))
	r.Invoke(`transferAsset`, users.TransferAsset, param.Struct(`data`, &users.GetTransaction{}))
	r.Invoke(`addAddress`, users.AddAddress, param.Struct(`data`, &users.Address{}))
	r.Invoke(`sendBalance`, users.TransferBalance, param.Struct(`data`, &users.SendBalance{}))

	// return the routes
	return chaincode
}

// Invoked when the chaincode is instantiated or upgraded
func invokeInit(c router.Context) (interface{}, error) {
	return owner.SetFromCreator(c)
}

// Execution start point
func main() {
	if err := shim.Start(New()); err != nil {
		fmt.Printf("Error starting Walletdemo chaincode: %s", err)
	}
}
