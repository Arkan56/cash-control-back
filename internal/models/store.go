package models

type Store struct {
	ID      int     `json:"id" db:"id"`
	Name    string  `json:"name" db:"name"`
	Balance float32 `json:"balance" db:"balance"`
}
