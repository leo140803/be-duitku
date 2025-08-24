package models

type Transaction struct {
    ID          string  `json:"id,omitempty"`
    UserID      string  `json:"user_id"`
    AccountID   string  `json:"account_id"`
    CategoryID  string  `json:"category_id,omitempty"`
    Date        string  `json:"date"`
    Description string  `json:"description,omitempty"`
    Amount      float64 `json:"amount"`
    Type        string  `json:"type"` // "INCOME" atau "EXPENSE"
    BalanceAfter float64 `json:"balance_after"`
    CreatedAt   string  `json:"created_at,omitempty"`
}
