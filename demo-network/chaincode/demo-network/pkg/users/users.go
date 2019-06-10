// Package users Related functions
package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/chaincode/demo-network/pkg/core/status"
	"github.com/chaincode/demo-network/pkg/core/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/s7techlab/cckit/router"
)

type ResponseAddAsset struct {
	ID      string `json:"_id"`
	Balance int64  `json:"balance"`
}

// GetUser create the user
func GetUser(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(User)

	// set the default values for the fields
	data.DocType = utils.DocTypeUser
	data.WalletBalance = 10000
	data.CreatedAt = time.Now().Format(time.RFC3339)

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	// check if address already exists or not
	queryString := fmt.Sprintf("{\"selector\":{\"address\":\"%s\",\"doc_type\":\"%s\"}}", data.Address, utils.DocTypeUser)
	address, user_id, err := utils.Get(c, queryString, fmt.Sprintf("Address already exists with the given address %s!", data.Address))

	//If address not found
	if address == nil {
		// return nil, err
		// get the stub to use it for query and save
		stub := c.Stub()

		// prepare the response body
		responseBody := UserResponse{ID: stub.GetTxID(), Address: data.Address, WalletBalance: data.WalletBalance, CreatedAt: data.CreatedAt}

		// Save the data and return the response
		return responseBody, c.State().Put(stub.GetTxID(), data)
	}

	userData := UserResponse{}
	err = json.Unmarshal(address, &userData)
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
	asset, _, err := utils.Get(c, queryString, "")
	if asset != nil {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("Symbol %s already exists!", data.Code))
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
	responseBody := ResponseAddAsset{ID: txID, Balance: user.WalletBalance}

	// Save the data and return the response
	return responseBody, c.State().Put(data.UserID, user)
}

// CheckAsset to check asset is available or not
func CheckAsset(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(CheckAssetStruct)

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
	asset, _, err := utils.Get(c, queryString, "")
	if asset != nil {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("Symbol %s already exists!", data.Code))
	}

	responseBody := utils.ResponseMessage{Message: "Symbol Available!"}

	// return the response
	return responseBody, nil
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
	queryRecevierString := fmt.Sprintf("{\"selector\":{\"address\":\"%s\",\"doc_type\":\"%s\"}}", data.To, utils.DocTypeUser)
	receiverData, receiverID, err5 := utils.Get(c, queryRecevierString, fmt.Sprintf("Receiver %s does not exist!", data.To))
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
	senderData, _, err6 := utils.Get(c, querySenderString, fmt.Sprintf("You account %s does not exist!", data.From))
	if err6 != nil {
		return nil, err6
	}
	sender := User{}
	err = json.Unmarshal(senderData, &sender)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	if sender.Address == data.To {
		return nil, status.ErrInternal.WithMessage(fmt.Sprintf("You can't transfer asset to yourself!"))
	}

	// check sender asset data
	queryString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\",\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", data.Code, data.From, utils.DocTypeAsset)
	senderAssetData, senderAssetKey, err2 := utils.Get(c, queryString, fmt.Sprintf("Symbol %s does not exist!", data.Code))
	if senderAssetData == nil {
		return nil, err2
	}
	senderAsset := Asset{}
	err = json.Unmarshal(senderAssetData, &senderAsset)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}
	if data.Quantity > senderAsset.Quantity {
		return nil, status.ErrInternal.WithMessage(fmt.Sprintf("Quantity should be less or equal to %d", senderAsset.Quantity))
	}
	stub := c.Stub()
	txID := stub.GetTxID()
	data.CreatedAt = time.Now().Format(time.RFC3339)
	// sender transactions
	var senderTransaction = Transaction{UserID: data.From, Type: 1, Code: data.Code, Quantity: data.Quantity, DocType: utils.DocTypeTransaction, CreatedAt: data.CreatedAt}
	err = c.State().Put(txID, senderTransaction)
	if err != nil {
		return nil, err
	}

	// receiver transactions
	var receiveTransaction = Transaction{UserID: receiverID, Type: 2, Code: data.Code, Quantity: data.Quantity, DocType: utils.DocTypeTransaction, CreatedAt: data.CreatedAt}
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
	queryReceiverDataString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\",\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", data.Code, receiverID, utils.DocTypeAsset)
	receiverAssetData, receiveAssetKey, _ := utils.Get(c, queryReceiverDataString, "")
	if receiverAssetData == nil {
		// add to receiver asset
		var receiveAsset = Asset{UserID: receiverID, Code: data.Code, Quantity: data.Quantity, DocType: utils.DocTypeAsset}
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
