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
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the user structure, with 6 properties.  Structure tags are used by encoding/json library
type User struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	WalletBalance int64  `json:"wallet_balance"`
	DocType       string `json:"doc_type"`
}

// Define the asset structure
type Asset struct {
	UserID   string `json:"user_id"`
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
	DocType  string `json:"doc_type"`
}

// Define the transactions structure
type Transaction struct {
	UserID   string `json:"user_id"`
	Type     int32  `json:"type"`
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
	DocType  string `json:"doc_type"`
}

// Define the user structure, with 6 properties.  Structure tags are used by encoding/json library
type UserResponse struct {
	ID            string `json:"_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	WalletBalance int64  `json:"wallet_balance"`
	DocType       string `json:"doc_type"`
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
	if function == "createUser" {
		return s.createUser(APIstub, args)
	} else if function == "getUser" {
		return s.getUser(APIstub, args)
	} else if function == "getUsers" {
		return s.getUsers(APIstub, args)
	} else if function == "getAssets" {
		return s.getAssets(APIstub, args)
	} else if function == "addAsset" {
		return s.addAsset(APIstub, args)
	} else if function == "transferAsset" {
		return s.transferAsset(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) createUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var user = User{Name: args[0], Email: args[1], Phone: args[2], Address: args[3], WalletBalance: 100, DocType: "users"}

	userAsBytes, _ := json.Marshal(user)
	txId := APIstub.GetTxID()
	APIstub.PutState(txId, userAsBytes)

	updatedUser := UserResponse{}
	err := json.Unmarshal(userAsBytes, &updatedUser)
	if err != nil {
		return shim.Error(err.Error())
	}
	updatedUser.ID = txId
	userAsBytes, _ = json.Marshal(updatedUser)

	return shim.Success(userAsBytes)
}

func (s *SmartContract) getUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userAsBytes, err := APIstub.GetState(args[0])
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}

	if userAsBytes != nil {
		return shim.Success(userAsBytes)
	}

	return shim.Error("user not found.")

}

func (s *SmartContract) getUsers(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"_id\":{\"$ne\":\"%s\"},\"doc_type\":\"%s\"}}", args[0], "users")
	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString("\"users\": [")
	aArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if aArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		userData := UserResponse{}
		err4 := json.Unmarshal(queryResponse.Value, &userData)
		if err4 != nil {
			return shim.Error(err4.Error())
		}

		userData.ID = queryResponse.Key
		userDataBytes, _ := json.Marshal(userData)

		buffer.WriteString(string(userDataBytes))
		aArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]}")
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) getAssets(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", args[0], "assets")
	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	queryTransactionsString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", args[0], "transactions")
	resultsIterator2, err := APIstub.GetQueryResult(queryTransactionsString)
	if err != nil {
		fmt.Println(err)
		return shim.Error(err.Error())
	}
	defer resultsIterator2.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString("\"assets\": [")
	aArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if aArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		aArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("],")

	buffer.WriteString("\"transactions\": [")

	bArrayMemberAlreadyWritten := false
	for resultsIterator2.HasNext() {
		queryResponse2, err := resultsIterator2.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse2.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]}")
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) addAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// check already exists
	queryString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\",\"doc_type\":\"%s\"}}", args[1], "assets")
	resultsIterator, err1 := APIstub.GetQueryResult(queryString)
	if err1 != nil {
		fmt.Println(err1)
		return shim.Error(err1.Error())
	}

	defer resultsIterator.Close()
	if resultsIterator.HasNext() {
		return shim.Error("This symbol already exists")
	}

	quantity, _ := strconv.Atoi(args[2])
	var add_asset = Asset{UserID: args[0], Code: args[1], Quantity: quantity, DocType: "assets"}
	assetAsBytes, _ := json.Marshal(add_asset)
	txId := APIstub.GetTxID()
	APIstub.PutState(txId, assetAsBytes)

	userAsBytes, _ := APIstub.GetState(args[0])
	user := User{}

	err := json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return shim.Error(err.Error())
	}

	user.WalletBalance = user.WalletBalance - 5

	userAsBytes, _ = json.Marshal(user)
	APIstub.PutState(args[0], userAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) transferAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// check already exists
	queryString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\",\"doc_type\":\"%s\"}}", args[2], "assets")
	resultsIterator, err1 := APIstub.GetQueryResult(queryString)
	if err1 != nil {
		fmt.Println(err1)
		return shim.Error(err1.Error())
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext() {
		return shim.Error("This symbol not exists")
	}

	txId := APIstub.GetTxID()
	quantity, _ := strconv.Atoi(args[3])
	// sender transactions
	var senderTransaction = Transaction{UserID: args[0], Type: 1, Code: args[2], Quantity: quantity, DocType: "transactions"}
	senderAssetAsBytes, _ := json.Marshal(senderTransaction)
	APIstub.PutState(txId, senderAssetAsBytes)

	// receiver transactions
	var receiveTransaction = Transaction{UserID: args[1], Type: 2, Code: args[2], Quantity: quantity, DocType: "transactions"}
	assetAsBytes, _ := json.Marshal(receiveTransaction)
	APIstub.PutState(txId+strconv.Itoa(1), assetAsBytes)

	// fetch the data and marshal it into struct
	queryResponse, err2 := resultsIterator.Next()
	if err2 != nil {
		return shim.Error(err2.Error())
	}
	assetData := Asset{}
	err := json.Unmarshal(queryResponse.Value, &assetData)
	if err != nil {
		return shim.Error(err.Error())
	}

	assetData.Quantity = assetData.Quantity - quantity
	// update sender asset data
	assetDataBytes, _ := json.Marshal(assetData)
	APIstub.PutState(queryResponse.Key, assetDataBytes)

	// check receiver data
	queryReceiverString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\",\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", args[2], args[1], "assets")
	resultsIterator2, err3 := APIstub.GetQueryResult(queryReceiverString)
	if err3 != nil {
		fmt.Println(err3)
		return shim.Error(err3.Error())
	}
	defer resultsIterator2.Close()

	// check exists or not
	if resultsIterator2.HasNext() {
		// fetch the data and marshal it into struct
		queryResponse2, err2 := resultsIterator2.Next()
		if err2 != nil {
			return shim.Error(err2.Error())
		}
		assetData := Asset{}
		err := json.Unmarshal(queryResponse2.Value, &assetData)
		if err != nil {
			return shim.Error(err.Error())
		}

		assetData.Quantity = assetData.Quantity + quantity
		// update sender asset data
		assetDataBytes, _ := json.Marshal(assetData)
		APIstub.PutState(queryResponse2.Key, assetDataBytes)
	} else {
		// add to receiver asset
		var receiveAsset = Asset{UserID: args[1], Code: args[2], Quantity: quantity, DocType: "assets"}
		assetAsBytes, _ := json.Marshal(receiveAsset)
		APIstub.PutState(txId+strconv.Itoa(3), assetAsBytes)
	}
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
