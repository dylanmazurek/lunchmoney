package models

type User struct {
	UserName    string  `json:"user_name"`
	UserEmail   string  `json:"user_email"`
	UserID      int     `json:"user_id"`
	AccountID   int     `json:"account_id"`
	BudgetName  *string `json:"budget_name"`
	ApiKeyLabel *string `json:"api_key_label"`
}
