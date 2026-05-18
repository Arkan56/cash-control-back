package models

import "time"

type Movement struct {
	ID               int64     `json:"id" db:"id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	Detail           string    `json:"detail" db:"detail"`
	Amount           int64     `json:"amount" db:"amount"`
	AmountCategoryID int32     `json:"amount_category_id" db:"amount_category_id"`
	VaultID          int64     `json:"vault_id" db:"vault_id"`
	UserID           int64     `json:"user_id" db:"user_id"`
}

type MovementWithUser struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Detail    string    `json:"detail"`
	Amount    int64     `json:"amount"`
	UserName  string    `json:"user_name"`
}
