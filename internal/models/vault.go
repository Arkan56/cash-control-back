package models

type Vault struct {
	ID       int64  `json:"id" db:"id"`
	StoreId  int64  `json:"store_id" db:"store_id"`
	Name     string `json:"name" db:"name"`
	Password string `json:"password" db:"password"`
	Balance  int64  `json:"balance" db:"balance"`
}

type CreateVaultRequest struct {
	StoreId  int64  `json:"store_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreatedVaultResponse struct {
	ID      int64  `json:"id"`
	StoreId int64  `json:"store_id"`
	Name    string `json:"name"`
	Balance int64  `json:"balance" db:"balance"`
}

type StoreVault struct {
	ID        int64  `json:"id"`
	StoreName string `json:"store_name"`
	Name      string `json:"name"`
}

type StoreVaultFull struct {
	ID        int64  `json:"id"`
	StoreName string `json:"store_name"`
	Name      string `json:"name"`
	Password  string
	Balance   int64 `json:"balance"`
}
