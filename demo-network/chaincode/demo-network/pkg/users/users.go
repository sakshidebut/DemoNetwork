// Package users Related functions
package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/chaincode/demo-network/pkg/core/status"
	"github.com/chaincode/demo-network/pkg/core/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/s7techlab/cckit/router"
)

// CreateUser create the user
func CreateUser(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(User)

	// set the default values for the fields
	data.DocType = utils.DocTypeUser
	data.WalletBalance = 10000

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	// check the user already exists or not
	queryString := fmt.Sprintf("{\"selector\":{\"email\":\"%s\",\"doc_type\":\"%s\"}}", data.Email, utils.DocTypeUser)
	alreadyExists, _, err := utils.Get(c, queryString, fmt.Sprintf("User already exists with email id %s!", data.Email))
	if alreadyExists != nil {
		return nil, err
	}

	// get the stub to use it for query and save
	stub := c.Stub()

	// prepare the response body
	responseBody := UserResponse{ID: stub.GetTxID(), Name: data.Name, Email: data.Email, Phone: data.Phone, Address: data.Address, WalletBalance: data.WalletBalance}

	// Save the data and return the response
	return responseBody, c.State().Put(stub.GetTxID(), data)
}

// GetUser get the user
func GetUser(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(UserId)

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	// check the user already exists or not
	queryString := fmt.Sprintf("{\"selector\":{\"email\":\"%s\",\"doc_type\":\"%s\"}}", data.ID, utils.DocTypeUser)
	user, user_id, err := utils.Get(c, queryString, fmt.Sprintf("User does not already exists with email id %s!", data.ID))
	if user == nil {
		return nil, err
	}

	userData := UserResponse{}
	err = json.Unmarshal(user, &userData)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}
	userData.ID = user_id

	userBytes, _ := json.Marshal(userData)

	//return the response
	return userBytes, nil
}

// GetUsers get the all users
func GetUsers(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(UserId)

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}
	stub := c.Stub()
	queryString := fmt.Sprintf("{\"selector\":{\"_id\":{\"$ne\":\"%s\"},\"doc_type\":\"%s\"}}", data.ID, utils.DocTypeUser)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		fmt.Println(err)
		return nil, status.ErrInternal.WithError(err)
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString("\"users\": [")
	aArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err2 := resultsIterator.Next()
		if err2 != nil {
			return nil, status.ErrInternal.WithError(err2)
		}

		// Add a comma before array members, suppress it for the first array member
		if aArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		userData := UserResponse{}
		err3 := json.Unmarshal(queryResponse.Value, &userData)
		if err3 != nil {
			return nil, status.ErrInternal.WithError(err3)
		}

		userData.ID = queryResponse.Key
		userDataBytes, _ := json.Marshal(userData)

		buffer.WriteString(string(userDataBytes))
		aArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]}")

	//return the response
	return buffer.Bytes(), nil
}

// GetAssets get the all Assets of user
func GetAssets(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(UserId)

	// Validate the inputed data
	err := data.Validate()
	if err != nil {

		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	stub := c.Stub()
	queryString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", data.ID, utils.DocTypeAsset)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		fmt.Println(err)
		return nil, status.ErrInternal.WithError(err)
	}
	defer resultsIterator.Close()

	queryTransactionsString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", data.ID, utils.DocTypeTransaction)
	resultsIterator2, err := stub.GetQueryResult(queryTransactionsString)
	if err != nil {
		fmt.Println(err)
		return nil, status.ErrInternal.WithError(err)
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
			return nil, status.ErrInternal.WithError(err)
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
			return nil, status.ErrInternal.WithError(err)
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse2.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]}")

	//return the response
	return buffer.Bytes(), nil
}

// AddAsset to add asset by user
func AddAsset(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(Asset)

	// set the default values for the fields
	data.DocType = utils.DocTypeAsset

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	// check already exists
	queryString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\",\"doc_type\":\"%s\"}}", data.Code, utils.DocTypeAsset)
	asset, _, err := utils.Get(c, queryString, fmt.Sprintf("Symbol %s already exists!", data.Code))
	if asset != nil {
		return nil, err
	}

	stub := c.Stub()
	txID := stub.GetTxID()

	err = c.State().Put(txID, data)
	if err != nil {
		return nil, err
	}

	userAsBytes, _ := stub.GetState(data.UserID)
	user := User{}

	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return nil, err
	}

	user.WalletBalance = user.WalletBalance - 5
	responseBody := utils.ResponseID{ID: txID}

	// Save the data and return the response
	return responseBody, c.State().Put(data.UserID, user)
}

// TransferAsset to transfer asset to another user
func TransferAsset(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(GetTransaction)

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	// check receiver data
	queryRecevierString := fmt.Sprintf("{\"selector\":{\"_id\":\"%s\",\"doc_type\":\"%s\"}}", data.To, utils.DocTypeUser)
	receiverData, _, err5 := utils.Get(c, queryRecevierString, fmt.Sprintf("Receiver %s does not exists!"))
	if err5 != nil {
		return nil, err5
	}
	receiver := User{}
	err = json.Unmarshal(receiverData, &receiver)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	// check sender data
	querySenderString := fmt.Sprintf("{\"selector\":{\"_id\":\"%s\",\"doc_type\":\"%s\"}}", data.From, utils.DocTypeUser)
	senderData, _, err6 := utils.Get(c, querySenderString, fmt.Sprintf("You account %s does not exists!"))
	if err6 != nil {
		return nil, err6
	}
	sender := User{}
	err = json.Unmarshal(senderData, &sender)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	// check sender asset data
	queryString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\",\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", data.Code, data.From, utils.DocTypeAsset)
	senderAssetData, senderAssetKey, err2 := utils.Get(c, queryString, fmt.Sprintf("Symbol %s does not exists!", data.Code))
	if senderAssetData == nil {
		return nil, err2
	}
	senderAsset := Asset{}
	err = json.Unmarshal(senderAssetData, &senderAsset)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}
	if data.Quantity > senderAsset.Quantity {
		return nil, status.ErrInternal.WithMessage(fmt.Sprintf("Quantity should be less or equal to %s", senderAsset.Quantity))
	}
	stub := c.Stub()
	txID := stub.GetTxID()

	// sender transactions
	var senderTransaction = Transaction{UserName: receiver.Name, UserID: data.From, Type: 1, Code: data.Code, Quantity: data.Quantity, DocType: utils.DocTypeTransaction}
	err = c.State().Put(txID, senderTransaction)
	if err != nil {
		return nil, err
	}

	// receiver transactions
	var receiveTransaction = Transaction{UserName: sender.Name, UserID: data.To, Type: 2, Code: data.Code, Quantity: data.Quantity, DocType: utils.DocTypeTransaction}
	err = c.State().Put(txID+strconv.Itoa(1), receiveTransaction)
	if err != nil {
		return nil, err
	}

	senderAsset.Quantity = senderAsset.Quantity - data.Quantity

	// update sender asset data
	err = c.State().Put(senderAssetKey, senderAsset)
	if err != nil {
		return nil, err
	}

	// check receiver asset data
	queryReceiverDataString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\",\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", data.Code, data.To, utils.DocTypeAsset)
	receiverAssetData, receiveAssetKey, _ := utils.Get(c, queryReceiverDataString, "")
	if receiverAssetData == nil {
		// add to receiver asset
		var receiveAsset = Asset{UserID: data.To, Code: data.Code, Quantity: data.Quantity, DocType: utils.DocTypeAsset}
		err = c.State().Put(txID+strconv.Itoa(3), receiveAsset)
		if err != nil {
			return nil, err
		}
	} else {
		receiverAsset := Asset{}
		err = json.Unmarshal(receiverAssetData, &receiverAsset)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		// update receiver asset
		receiverAsset.Quantity = receiverAsset.Quantity + data.Quantity
		err = c.State().Put(receiveAssetKey, receiverAsset)
		if err != nil {
			return nil, err
		}
	}

	responseBody := utils.ResponseID{ID: txID}

	// return the response
	return responseBody, nil
}
