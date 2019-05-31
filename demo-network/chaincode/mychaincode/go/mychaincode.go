/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the user structure, with 6 properties.  Structure tags are used by encoding/json library
type User struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone string `json:"phone"`
	Currency  string `json:"currency"`
	WalletBalance  string `json:"wallet_balance"`
	DocType  string `json:"doc_type"`
}

// Define the issue-token structure, with 5 properties.  Structure tags are used by encoding/json library
type IssueToken struct {
	UserID   string `json:"user_id"`
	Code  string `json:"code"`
	Currency  string `json:"currency"`
	Quantity string `json:"quantity"`
	DocType  string `json:"doc_type"`
}

/*
 * The Init method is called when the Smart Contract is instantiated by the blockchain network
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "getWalletBalance" {
		return s.getWalletBalance(APIstub, args)
	} else if function == "createUser" {
		return s.createUser(APIstub, args)
	} else if function == "issueToken" {
		return s.issueToken(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) getWalletBalance(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(userAsBytes)
}

func (s *SmartContract) createUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var user = User{Name: args[0], Email: args[1], Phone: args[2], Currency: "INR", WalletBalance: "100", DocType: "users"}

	userAsBytes, _ := json.Marshal(user)
	txId := APIstub.GetTxID()
	APIstub.PutState(txId, userAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) issueToken(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var issue_token = IssueToken{UserID: args[0], Code: args[1], Currency: "INR", Quantity: args[2], DocType: "issue_tokens"}

	issueTokenAsBytes, _ := json.Marshal(issue_token)
	txId := APIstub.GetTxID()
	APIstub.PutState(txId, issueTokenAsBytes)

	userAsBytes, _ := APIstub.GetState(args[0])
	user := User{}

	json.Unmarshal(userAsBytes, &user)
	// Wval, _ = strconv.Atoi(user.WalletBalance)
	// Uval, _ = strconv.Atoi(Wval - 2)
	// user.WalletBalance = Uval

	userAsBytes, _ = json.Marshal(user)
	APIstub.PutState(args[0], userAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
