// Package  users Validations
package users

import (
	"github.com/chaincode/demo-network/pkg/core/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Validate Validates the User Structure
func (data User) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Address, validation.Required.Error(utils.AddressRequired), validation.NotNil.Error(utils.AddressRequired)),
	)
}

// Validate Validates the Address Structure
func (data Address) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Label, validation.Required.Error(utils.AddressRequired), validation.NotNil.Error(utils.AddressRequired)),
		validation.Field(&data.Value, validation.Required.Error(utils.AddressRequired), validation.NotNil.Error(utils.AddressRequired)),
	)
}

// Validate Validates the Asset Structure
func (data Asset) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.UserID, validation.Required.Error(utils.UserIDRequired), validation.NotNil.Error(utils.UserIDRequired)),
		validation.Field(&data.Code, validation.Required.Error(utils.CodeRequired), validation.NotNil.Error(utils.CodeRequired)),
		validation.Field(&data.Quantity, validation.Required.Error(utils.QuantityRequired), validation.NotNil.Error(utils.QuantityRequired)),
	)
}

// Validate Validates the CheckAssetStruct Structure
func (data CheckAssetStruct) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Code, validation.Required.Error(utils.CodeRequired), validation.NotNil.Error(utils.CodeRequired)),
	)
}

// Validate Validates the UserId Structure
func (data UserId) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.ID, validation.Required.Error(utils.IDRequired), validation.NotNil.Error(utils.IDRequired)),
	)
}

// Validate Validates the UserId Structure
func (data GetTransaction) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.From, validation.Required.Error(utils.IDRequired), validation.NotNil.Error(utils.IDRequired)),
		validation.Field(&data.To, validation.Required.Error(utils.IDRequired), validation.NotNil.Error(utils.IDRequired)),
		validation.Field(&data.Code, validation.Required.Error(utils.CodeRequired), validation.NotNil.Error(utils.CodeRequired)),
		validation.Field(&data.Quantity, validation.Required.Error(utils.QuantityRequired), validation.NotNil.Error(utils.QuantityRequired)),
	)
}
