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
	Symbol        string    `json:"symbol"`
	DocType       string    `json:"doc_type"`
	CreatedAt     string    `json:"created_at"`
	UserAddresses []Address `json:"user_addresses"`
	Identity      string    `json:"identity"`
	Secret        string    `json:"secret"`
}

// Define the asset structure
type Asset struct {
	UserID   string `json:"user_id"`
	Label    string `json:"label"`
	Code     string `json:"code"`
	Quantity int64  `json:"quantity"`
	DocType  string `json:"doc_type"`
}

// Define the CheckAssetStruct structure
type CheckAssetStruct struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// Define the transactions structure
type Transaction struct {
	UserID           string `json:"user_id"`
	TxnType          string `json:"txn_type"`
	Type             int32  `json:"type"`
	Code             string `json:"code"`
	AssetLabel       string `json:"asset_label"`
	Quantity         int64  `json:"quantity"`
	AddressValue     string `json:"address_value"`
	LabelValue       string `json:"label_value"`
	AddressBookLabel string `json:"address_book_label"`
	DocType          string `json:"doc_type"`
	CreatedAt        string `json:"created_at"`
}

// Define the user structure, with 6 properties.  Structure tags are used by encoding/json library
type NewUserResponse struct {
	ID            string    `json:"_id"`
	Address       string    `json:"address"`
	WalletBalance int64     `json:"wallet_balance"`
	Symbol        string    `json:"symbol"`
	CreatedAt     string    `json:"created_at"`
	UserAddresses []Address `json:"user_addresses"`
	Identity      string    `json:"identity"`
	Secret        string    `json:"secret"`
}

// Define the user structure, with 6 properties.  Structure tags are used by encoding/json library
type UserResponse struct {
	ID            string    `json:"_id"`
	Address       string    `json:"address"`
	WalletBalance int64     `json:"wallet_balance"`
	Symbol        string    `json:"symbol"`
	CreatedAt     string    `json:"created_at"`
	UserAddresses []Address `json:"user_addresses"`
	Identity      string    `json:"identity"`
}

// Define the UserId structure
type UserId struct {
	ID string `json:"id"`
}

// Define the UserSecret structure
type UserSecret struct {
	Secret string `json:"secret"`
}

// Define the GetTransactions structure
type GetTransaction struct {
	From      string `json:"from_id"`
	To        string `json:"to_id"`
	Code      string `json:"code"`
	Quantity  int64  `json:"quantity"`
	Label     string `json:"label"`
	DocType   string `json:"doc_type"`
	CreatedAt string `json:"created_at"`
}

type ResponseAddAsset struct {
	ID      string `json:"_id"`
	Balance int64  `json:"balance"`
	Symbol  string `json:"symbol"`
}

// Define the transactions structure
type TransactionResponse struct {
	ID               string `json:"_id"`
	UserID           string `json:"user_id"`
	TxnType          string `json:"txn_type"`
	Type             int32  `json:"type"`
	Code             string `json:"code"`
	AssetLabel       string `json:"asset_label"`
	Quantity         int64  `json:"quantity"`
	AddressValue     string `json:"address_value"`
	LabelValue       string `json:"label_value"`
	AddressBookLabel string `json:"address_book_label"`
	DocType          string `json:"doc_type"`
	CreatedAt        string `json:"created_at"`
}

// Define the SendBalance structure
type SendBalance struct {
	From     string `json:"from_id"`
	To       string `json:"to_id"`
	Quantity int64  `json:"quantity"`
	Label    string `json:"label"`
}

// Define the AddressBook structure
type AddressBook struct {
	UserID  string `json:"user_id"`
	Address string `json:"address"`
	Label   string `json:"label"`
	DocType string `json:"doc_type"`
}
