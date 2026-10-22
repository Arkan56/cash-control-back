package models

type AmountCategory struct {
	ID          uint   `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
}
