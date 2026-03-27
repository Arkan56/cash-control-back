package models

import "time"

type Transaction struct {
	ID               uint
	Date             time.Time
	Datail           string
	Amount           float32
	IdAmountCategory uint
	IdVault          uint
}
