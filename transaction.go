package mango

// Txn represents a transaction in the Manifold system.
type Txn struct {
	Id          string  `json:"id"`
	CreatedTime int64   `json:"createdTime"`
	FromId      string  `json:"fromId"`
	FromType    string  `json:"fromType"`
	ToId        string  `json:"toId"`
	ToType      string  `json:"toType"`
	Amount      float64 `json:"amount"`
	Token       string  `json:"token"`
	Category    string  `json:"category"`
	Description string  `json:"description,omitempty"`
}

// GetTransactionsRequest represents the parameters for fetching transactions.
type GetTransactionsRequest struct {
	Token    string `json:"token,omitempty"`
	Offset   int64  `json:"offset,omitempty"`
	Limit    int64  `json:"limit,omitempty"`
	Before   int64  `json:"before,omitempty"`
	After    int64  `json:"after,omitempty"`
	ToId     string `json:"toId,omitempty"`
	FromId   string `json:"fromId,omitempty"`
	Category string `json:"category,omitempty"`
}

// SendManagramRequest represents the parameters for sending mana to users.
type SendManagramRequest struct {
	ToIds   []string `json:"toIds"`
	Amount  float64  `json:"amount"`
	Message string   `json:"message,omitempty"`
}
