// Package users Related functions
package users

type Address struct {
	UserID string `json:"user_id"`
	Label  string `json:"label"`
	Value  string `json:"value"`
}

// Define the user structure, with 6 properties.  Structure tags are used by encoding/json library
type User struct {
	Address       string    `json:"address"`
	WalletBalance int64     `json:"wallet_balance"`
	DocType       string    `json:"doc_type"`
	CreatedAt     string    `json:"created_at"`
	UserAddresses []Address `json:"user_addresses"`
}

// Define the asset structure
type Asset struct {
	UserID   string `json:"user_id"`
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
	DocType  string `json:"doc_type"`
}

// Define the CheckAssetStruct structure
type CheckAssetStruct struct {
	Code string `json:"code"`
}

// Define the transactions structure
type Transaction struct {
	UserID       string `json:"user_id"`
	Type         int32  `json:"type"`
	Code         string `json:"code"`
	Quantity     int    `json:"quantity"`
	AddressValue string `json:"address_value"`
	DocType      string `json:"doc_type"`
	CreatedAt    string `json:"created_at"`
}

// Define the user structure, with 6 properties.  Structure tags are used by encoding/json library
type UserResponse struct {
	ID            string    `json:"_id"`
	Address       string    `json:"address"`
	WalletBalance int64     `json:"wallet_balance"`
	CreatedAt     string    `json:"created_at"`
	UserAddresses []Address `json:"user_addresses"`
}

// Define the UserId structure
type UserId struct {
	ID string `json:"id"`
}

// Define the GetTransactions structure
type GetTransaction struct {
	From      string `json:"from_id"`
	To        string `json:"to_id"`
	Code      string `json:"code"`
	Quantity  int    `json:"quantity"`
	DocType   string `json:"doc_type"`
	CreatedAt string `json:"created_at"`
}

type ResponseAddAsset struct {
	ID      string `json:"_id"`
	Balance int64  `json:"balance"`
}
