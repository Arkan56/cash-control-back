package models

type Vault struct {
	ID      int64  `json:"id" db:"id"`
	StoreId int64  `json:"store_id" db:"store_id"`
	Name    string `json:"name" db:"name"`
	Balance int64  `json:"balance" db:"balance"`
}

type CreateVaultRequest struct {
	StoreId int64  `json:"store_id" binding:"required"`
	Name    string `json:"name" binding:"required"`
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

type StoreIdVault struct {
	ID      int64  `json:"id"`
	StoreId int64  `json:"store_id"`
	Name    string `json:"name"`
}

type StoreVaultFull struct {
	ID        int64  `json:"id"`
	StoreName string `json:"store_name"`
	Name      string `json:"name"`
	Balance   int64  `json:"balance"`
}

type VaultAccessResponse struct {
	VaultID int64 `json:"vault_id"`
}
