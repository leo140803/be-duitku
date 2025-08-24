package models

type Account struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	Name           string  `json:"name"`
	InitialBalance float64 `json:"initial_balance"`
	CreatedAt      string  `json:"created_at"`
}

// struct khusus untuk update saldo
type UpdateAccountBalance struct {
	InitialBalance float64 `json:"initial_balance"`
}
