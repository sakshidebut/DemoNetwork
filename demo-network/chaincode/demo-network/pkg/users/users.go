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

// CreateUser create the user
func CreateUser(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(User)

	// set the default values for the fields
	data.DocType = utils.DocTypeUser
	data.WalletBalance = 10000
	data.Symbol = utils.WalletCoinSymbol
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
	queryString := fmt.Sprintf("{\"selector\": {\"user_addresses\": {\"$elemMatch\": {\"value\": \"%s\"}},\"doc_type\":\"%s\"}}", data.Address, utils.DocTypeUser)
	address, userID, err := utils.Get(c, queryString, fmt.Sprintf("Address already exists with the given address %s!", data.Address))

	//If address not found
	if address == nil {
		// get the stub to use it for query and save
		stub := c.Stub()

		var addresses []Address
		address1 := Address{UserID: stub.GetTxID(), Label: "Original", Value: data.Address}
		addresses = append(addresses, address1)
		data.UserAddresses = addresses

		// prepare the response body
		responseBody := NewUserResponse{ID: stub.GetTxID(), Address: data.Address, WalletBalance: data.WalletBalance, Symbol: data.Symbol, CreatedAt: data.CreatedAt, UserAddresses: addresses, Identity: data.Identity, Secret: data.Secret + data.Identity}

		// Save the data and return the response
		return responseBody, c.State().Put(stub.GetTxID(), data)
	}

	userData := UserResponse{}
	err = json.Unmarshal(address, &userData)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}
	userData.ID = userID

	userBytes, _ := json.Marshal(userData)

	//return the response
	return userBytes, nil
}

// GetUser fetch the details of user
func GetUser(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(UserSecret)

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	secret := data.Secret[0:300]

	// check if user already exists or not
	queryString := fmt.Sprintf("{\"selector\":{\"secret\":\"%s\",\"doc_type\":\"%s\"}}", secret, utils.DocTypeUser)
	userResult, userID, err := utils.Get(c, queryString, fmt.Sprintf("User found!"))

	if userResult == nil {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("User does not exist in this system!"))
	}

	userData := UserResponse{}
	err = json.Unmarshal(userResult, &userData)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}
	userData.ID = userID

	userBytes, _ := json.Marshal(userData)

	//return the response
	return userBytes, nil
}

// AddAddress : Function to add the multiple addresses of user
func AddAddress(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(Address)

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	// check if address already exists or not
	queryString := fmt.Sprintf("{\"selector\": {\"user_addresses\": {\"$elemMatch\": {\"value\": \"%s\"}},\"doc_type\":\"%s\"}}", data.Value, utils.DocTypeUser)
	userResult, _, err := utils.Get(c, queryString, fmt.Sprintf("User already exists with the given address %s!", data.Value))

	if userResult != nil {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("This address %s already exists in the system!", data.Value))
	}

	// check if address already exists or not
	labelQueryString := fmt.Sprintf("{\"selector\": {\"user_addresses\": {\"$elemMatch\": {\"label\": \"%s\"}},\"_id\":{\"$ne\":\"%s\"},\"doc_type\":\"%s\"}}", data.Label, data.UserID, utils.DocTypeUser)
	labelResult, _, err := utils.Get(c, labelQueryString, fmt.Sprintf("User already exists with the given label %s!", data.Label))

	if labelResult != nil {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("This label %s has already been taken!", data.Label))
	}

	address1 := Address{UserID: data.UserID, Label: data.Label, Value: data.Value}
	stub := c.Stub()
	userAsBytes, _ := stub.GetState(data.UserID)
	user := User{}

	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return nil, err
	}

	user.UserAddresses = append(user.UserAddresses, address1)
	// prepare the response body
	responseBody := UserResponse{ID: data.UserID, Address: user.Address, WalletBalance: user.WalletBalance, Symbol: user.Symbol, CreatedAt: user.CreatedAt, UserAddresses: user.UserAddresses, Identity: user.Identity}
	// Save the data and return the response
	return responseBody, c.State().Put(data.UserID, user)
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
	queryUserString := fmt.Sprintf("{\"selector\":{\"_id\":\"%s\",\"doc_type\":\"%s\"}}", data.ID, utils.DocTypeUser)
	userData, _, err1 := utils.Get(c, queryUserString, fmt.Sprintf("User %s does not exist!", data.ID))
	if err1 != nil {
		return nil, err1
	}

	user := User{}
	err = json.Unmarshal(userData, &user)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}
	responseBody := ResponseAddAsset{ID: data.ID, Balance: user.WalletBalance, Symbol: user.Symbol}
	resBytes, _ := json.Marshal(responseBody)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	queryString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"doc_type\":\"%s\"}}", data.ID, utils.DocTypeAsset)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		fmt.Println(err)
		return nil, status.ErrInternal.WithError(err)
	}
	defer resultsIterator.Close()

	queryTransactionsString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"doc_type\":\"%s\"},\"sort\":[{\"created_at\":\"desc\"}]}", data.ID, utils.DocTypeTransaction)
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
	buffer.WriteString("\"wallet_balance\": ")
	buffer.WriteString(string(resBytes))
	buffer.WriteString(",")
	buffer.WriteString("\"transactions\": [")

	bArrayMemberAlreadyWritten := false
	for resultsIterator2.HasNext() {
		queryResponse2, err := resultsIterator2.Next()
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}

		txnData := TransactionResponse{}
		err3 := json.Unmarshal(queryResponse2.Value, &txnData)
		if err3 != nil {
			return nil, status.ErrInternal.WithError(err3)
		}

		txnData.ID = queryResponse2.Key
		txnDataBytes, _ := json.Marshal(txnData)

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(txnDataBytes))
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

	stub := c.Stub()
	txID := stub.GetTxID()
	userAsBytes, _ := stub.GetState(data.UserID)
	user := User{}

	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return nil, err
	}

	if user.WalletBalance < utils.AddAssetFee {
		return nil, status.ErrInternal.WithMessage(fmt.Sprintf("You don't have enough coins to purchase this asset."))
	}

	// check asset code already exists
	queryString := fmt.Sprintf("{\"selector\":{\"code\":\"%s\",\"doc_type\":\"%s\"}}", data.Code, utils.DocTypeAsset)
	asset, _, err := utils.Get(c, queryString, "")
	if asset != nil {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("Symbol %s already exists!", data.Code))
	}

	// check asset label already exists
	queryString1 := fmt.Sprintf("{\"selector\":{\"label\":\"%s\",\"doc_type\":\"%s\"}}", data.Label, utils.DocTypeAsset)
	assetLabel, _, err := utils.Get(c, queryString1, "")
	if assetLabel != nil {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("Name %s already exists!", data.Label))
	}

	err = c.State().Put(txID, data)
	if err != nil {
		return nil, err
	}

	user.WalletBalance = user.WalletBalance - utils.AddAssetFee

	createdAt := time.Now().Format(time.RFC3339)
	// add asset transaction
	var addAssetTransaction = Transaction{UserID: data.UserID, Type: utils.Send, Code: utils.WalletCoinSymbol, AssetLabel: data.Label, Quantity: utils.AddAssetFee, DocType: utils.DocTypeTransaction, CreatedAt: createdAt, AddressValue: "", LabelValue: "", AddressBookLabel: "Original", TxnType: utils.AssetCreatedTxn}
	err = c.State().Put(txID+strconv.Itoa(1), addAssetTransaction)
	if err != nil {
		return nil, err
	}

	responseBody := ResponseAddAsset{ID: txID, Balance: user.WalletBalance, Symbol: user.Symbol}

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

	// check already exists
	queryString1 := fmt.Sprintf("{\"selector\":{\"label\":\"%s\",\"doc_type\":\"%s\"}}", data.Label, utils.DocTypeAsset)
	asset1, _, err := utils.Get(c, queryString1, "")
	if asset1 != nil {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("Name %s already exists!", data.Label))
	}

	responseBody := utils.ResponseMessage{Message: "Both name and symbol are available."}

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
	queryRecevierString := fmt.Sprintf("{\"selector\": {\"user_addresses\": {\"$elemMatch\": {\"value\": \"%s\"}},\"doc_type\":\"%s\"}}", data.To, utils.DocTypeUser)
	receiverData, receiverID, err5 := utils.Get(c, queryRecevierString, fmt.Sprintf("Receiver %s does not exist!", data.To))
	if err5 != nil {
		return nil, err5
	}

	receiver := User{}
	err = json.Unmarshal(receiverData, &receiver)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}

	var receiverOwnLabel string
	for i := range receiver.UserAddresses {
		if receiver.UserAddresses[i].Value == data.To {
			receiverOwnLabel = receiver.UserAddresses[i].Label
		}
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

	if sender.WalletBalance < utils.TransferAssetFee {
		return nil, status.ErrInternal.WithMessage(fmt.Sprintf("You don't have enough coins to transfer the asset."))
	}

	for i := range sender.UserAddresses {
		if sender.UserAddresses[i].Value == data.To {
			return nil, status.ErrInternal.WithMessage(fmt.Sprintf("You can't transfer asset to yourself!"))
		}
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

	var receiverLabel, senderLabel string
	// check label of receiver in sender's address book
	receiverLabelString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"address\":\"%s\",\"label\":\"%s\",\"doc_type\":\"%s\"}}", data.From, data.To, data.Label, utils.DocTypeAddressBook)
	receiverLabelData, _, err6 := utils.Get(c, receiverLabelString, fmt.Sprintf("Label of receiver does not exist!"))

	//If label does not exist in address book then save it into db
	if receiverLabelData == nil {
		// check if label is unique
		checkUniqueString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"label\":\"%s\",\"doc_type\":\"%s\"}}", data.From, data.Label, utils.DocTypeAddressBook)
		uniqueLabelData, _, err := utils.Get(c, checkUniqueString, fmt.Sprintf("This label already exists!"))
		if uniqueLabelData != nil {
			return nil, status.ErrInternal.WithMessage(fmt.Sprintf("This label already exists!"))
		}

		labelTxn := AddressBook{UserID: data.From, Address: data.To, Label: data.Label, DocType: utils.DocTypeAddressBook}
		receiverLabel = data.Label
		// Save the data
		err = c.State().Put(txID, labelTxn)
		if err != nil {
			return nil, err
		}
	} else {

		addressLabel := AddressBook{}
		err = json.Unmarshal(receiverLabelData, &addressLabel)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		receiverLabel = addressLabel.Label
	}

	// check label of sender in receiver's address book
	senderLabelString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"address\":\"%s\",\"doc_type\":\"%s\"}}", receiver.UserAddresses[0].UserID, sender.Address, utils.DocTypeAddressBook)
	senderLabelData, _, err6 := utils.Get(c, senderLabelString, fmt.Sprintf("Label of sender does not exist!"))

	//If label does not exist in address book
	if senderLabelData == nil {
		senderLabel = "N/A"
	} else {

		addressLabel1 := AddressBook{}
		err = json.Unmarshal(senderLabelData, &addressLabel1)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		senderLabel = addressLabel1.Label
	}

	// sender transactions
	var senderTransaction = Transaction{UserID: data.From, Type: utils.Send, Code: data.Code, AssetLabel: senderAsset.Label, Quantity: data.Quantity, DocType: utils.DocTypeTransaction, CreatedAt: data.CreatedAt, AddressValue: data.To, LabelValue: receiverOwnLabel, AddressBookLabel: receiverLabel, TxnType: utils.AssetTxnType}
	err = c.State().Put(txID+strconv.Itoa(1), senderTransaction)
	if err != nil {
		return nil, err
	}

	// receiver transactions
	var receiveTransaction = Transaction{UserID: receiverID, Type: utils.Receive, Code: data.Code, AssetLabel: senderAsset.Label, Quantity: data.Quantity, DocType: utils.DocTypeTransaction, CreatedAt: data.CreatedAt, AddressValue: sender.Address, LabelValue: "Original", AddressBookLabel: senderLabel, TxnType: utils.AssetTxnType}
	err = c.State().Put(txID+strconv.Itoa(2), receiveTransaction)
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
		var receiveAsset = Asset{UserID: receiverID, Code: data.Code, Label: senderAsset.Label, Quantity: data.Quantity, DocType: utils.DocTypeAsset}
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

	sender.WalletBalance = sender.WalletBalance - utils.TransferAssetFee

	// Transfer asset transaction
	var transferAssetTransaction = Transaction{UserID: data.From, Type: utils.Send, Code: utils.WalletCoinSymbol, AssetLabel: senderAsset.Label, Quantity: utils.TransferAssetFee, DocType: utils.DocTypeTransaction, CreatedAt: data.CreatedAt, AddressValue: data.To, LabelValue: receiverOwnLabel, AddressBookLabel: receiverLabel, TxnType: utils.AssetTransferredTxn}
	err = c.State().Put(txID+strconv.Itoa(4), transferAssetTransaction)
	if err != nil {
		return nil, err
	}

	responseBody := ResponseAddAsset{ID: txID, Balance: sender.WalletBalance, Symbol: sender.Symbol}

	// Save the data and return the response
	return responseBody, c.State().Put(data.From, sender)
}

// TransferBalance to transfer asset to another user
func TransferBalance(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(SendBalance)

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	// check receiver data
	queryRecevierString := fmt.Sprintf("{\"selector\": {\"user_addresses\": {\"$elemMatch\": {\"value\": \"%s\"}},\"doc_type\":\"%s\"}}", data.To, utils.DocTypeUser)
	receiverData, _, err5 := utils.Get(c, queryRecevierString, fmt.Sprintf("Receiver %s does not exist!", data.To))
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

	for i := range sender.UserAddresses {
		if sender.UserAddresses[i].Value == data.To {
			return nil, status.ErrInternal.WithMessage(fmt.Sprintf("You can't transfer coins to yourself!"))
		}
	}

	if data.Quantity > sender.WalletBalance {
		return nil, status.ErrInternal.WithMessage(fmt.Sprintf("Quantity should be less or equal to %d", sender.WalletBalance))
	}

	stub := c.Stub()
	txID := stub.GetTxID()

	var receiverLabel, senderLabel string
	// check label of receiver in sender's address book
	receiverLabelString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"address\":\"%s\",\"label\":\"%s\",\"doc_type\":\"%s\"}}", data.From, data.To, data.Label, utils.DocTypeAddressBook)
	receiverLabelData, _, err6 := utils.Get(c, receiverLabelString, fmt.Sprintf("Label of receiver does not exist!"))

	//If label does not exist in address book then save it into db
	if receiverLabelData == nil {
		// check if label is unique
		checkUniqueString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"label\":\"%s\",\"doc_type\":\"%s\"}}", data.From, data.Label, utils.DocTypeAddressBook)
		uniqueLabelData, _, err := utils.Get(c, checkUniqueString, fmt.Sprintf("This label already exists!"))
		if uniqueLabelData != nil {
			return nil, status.ErrInternal.WithMessage(fmt.Sprintf("This label already exists!"))
		}

		labelTxn := AddressBook{UserID: data.From, Address: data.To, Label: data.Label, DocType: utils.DocTypeAddressBook}
		receiverLabel = data.Label
		// Save the data
		err = c.State().Put(txID, labelTxn)
		if err != nil {
			return nil, err
		}
	} else {

		addressLabel := AddressBook{}
		err = json.Unmarshal(receiverLabelData, &addressLabel)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		receiverLabel = addressLabel.Label
	}

	// check label of sender in receiver's address book
	senderLabelString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"address\":\"%s\",\"doc_type\":\"%s\"}}", receiver.UserAddresses[0].UserID, sender.Address, utils.DocTypeAddressBook)
	senderLabelData, _, err6 := utils.Get(c, senderLabelString, fmt.Sprintf("Label of sender does not exist!"))

	//If label does not exist in address book
	if senderLabelData == nil {
		senderLabel = "N/A"
	} else {

		addressLabel1 := AddressBook{}
		err = json.Unmarshal(senderLabelData, &addressLabel1)
		if err != nil {
			return nil, status.ErrInternal.WithError(err)
		}
		senderLabel = addressLabel1.Label
	}

	createdAt := time.Now().Format(time.RFC3339)
	// sender transactions
	var senderTransaction = Transaction{UserID: data.From, Type: utils.Send, Code: utils.WalletCoinSymbol, Quantity: data.Quantity, DocType: utils.DocTypeTransaction, CreatedAt: createdAt, AddressValue: data.To, LabelValue: "", AddressBookLabel: receiverLabel, TxnType: utils.CoinTxnType}
	err = c.State().Put(txID+strconv.Itoa(1), senderTransaction)
	if err != nil {
		return nil, err
	}

	// receiver transactions
	var receiveTransaction = Transaction{UserID: receiver.UserAddresses[0].UserID, Type: utils.Receive, Code: utils.WalletCoinSymbol, Quantity: data.Quantity, DocType: utils.DocTypeTransaction, CreatedAt: createdAt, AddressValue: sender.Address, LabelValue: "Original", AddressBookLabel: senderLabel, TxnType: utils.CoinTxnType}
	err = c.State().Put(txID+strconv.Itoa(2), receiveTransaction)
	if err != nil {
		return nil, err
	}

	// update sender wallet
	sender.WalletBalance = sender.WalletBalance - data.Quantity
	err = c.State().Put(data.From, sender)
	if err != nil {
		return nil, err
	}

	// update receiver wallet
	receiver.WalletBalance = receiver.WalletBalance + data.Quantity
	err = c.State().Put(receiver.UserAddresses[0].UserID, receiver)
	if err != nil {
		return nil, err
	}

	responseBody := ResponseAddAsset{ID: data.From, Balance: sender.WalletBalance, Symbol: sender.Symbol}
	// return the response
	return responseBody, nil
}

// GetAddressBookLabel to fetch address or label depending upon input
func GetAddressBookLabel(c router.Context) (interface{}, error) {
	// get the data from the request and parse it as structure
	data := c.Param(`data`).(AddressBook)

	// Validate the inputed data
	err := data.Validate()
	if err != nil {
		if _, ok := err.(validation.InternalError); ok {
			return nil, err
		}
		return nil, status.ErrStatusUnprocessableEntity.WithValidationError(err.(validation.Errors))
	}

	var queryLabelString string
	//If label is empty and address is not empty
	if (data.Label == "" && data.Address != "") || (data.Label != "" && data.Address != "") {
		// check label in address book corresponding to address
		queryLabelString = fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"address\":\"%s\",\"doc_type\":\"%s\"}}", data.UserID, data.Address, utils.DocTypeAddressBook)
	} else {
		// check address in address book corresponding to label
		queryLabelString = fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\",\"label\":\"%s\",\"doc_type\":\"%s\"}}", data.UserID, data.Label, utils.DocTypeAddressBook)
	}

	LabelData, _, err := utils.Get(c, queryLabelString, fmt.Sprintf("Record does not exist in your address book."))

	if LabelData == nil {
		return nil, status.ErrInternal.WithMessage(fmt.Sprintf("Label does not exist for this address."))
	}

	addressLabel := AddressBook{}
	err = json.Unmarshal(LabelData, &addressLabel)
	if err != nil {
		return nil, status.ErrInternal.WithError(err)
	}
	// return the response
	return addressLabel, nil
}
