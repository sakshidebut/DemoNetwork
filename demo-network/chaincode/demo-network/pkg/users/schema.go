// Package users Related functions
package users

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
	UserName string `json:"user_name"`
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
}

// Define the UserId structure
type UserId struct {
	ID string `json:"id"`
}

// Define the GetTransactions structure
type GetTransaction struct {
	From     string `json:"from_id"`
	To       string `json:"to_id"`
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
	DocType  string `json:"doc_type"`
}
