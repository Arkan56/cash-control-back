package models

type Store struct {
	ID      int     `json:"id" db:"id"`
	Name    string  `json:"name" db:"name"`
	Balance float32 `json:"balance" db:"balance"`
}

type StoreAccessResponse struct {
	StoreID int64 `json:"store_id"`
}

type StoreUpdateResponse struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type StoreDetail struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float32 `json:"balance"`
	Vaults  []Vault `json:"vaults"`
}
