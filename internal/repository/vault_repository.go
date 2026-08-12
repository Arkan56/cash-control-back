package repository

import (
	"cash-control/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateVault(pool *pgxpool.Pool, vault *models.CreateVaultRequest) (*models.CreatedVaultResponse, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	INSERT INTO vaults (
	store_id,
	name
	)
	VALUES($1, $2)
	RETURNING id, store_id, name, balance
	`

	var createdVault models.CreatedVaultResponse

	var err error = pool.QueryRow(ctx, query, vault.StoreId, vault.Name).Scan(
		&createdVault.ID,
		&createdVault.StoreId,
		&createdVault.Name,
		&createdVault.Balance,
	)

	if err != nil {
		return nil, err
	}

	return &createdVault, nil
}

func GetVaultsByStoreId(pool *pgxpool.Pool, StoreId int64) ([]models.StoreVault, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT vaults.id, stores.name, vaults.name
		FROM vaults
		JOIN stores ON vaults.store_id = stores.id
		WHERE vaults.store_id = $1
	`
	var rows, err = pool.Query(ctx, query, StoreId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var vaults []models.StoreVault = []models.StoreVault{}

	for rows.Next() {
		var v models.StoreVault

		err = rows.Scan(
			&v.ID,
			&v.StoreName,
			&v.Name,
		)

		if err != nil {
			return nil, err
		}

		vaults = append(vaults, v)
	}
	return vaults, nil
}

func GetVaultById(pool *pgxpool.Pool, vaultId int64) (*models.StoreVaultFull, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT vaults.id, stores.name, vaults.name, vaults.balance
		FROM vaults
		JOIN stores ON vaults.store_id = stores.id
		WHERE vaults.id = $1
	`
	var vault models.StoreVaultFull
	err := pool.QueryRow(ctx, query, vaultId).Scan(
		&vault.ID,
		&vault.StoreName,
		&vault.Name,
		&vault.Balance,
	)

	if err != nil {
		return nil, err
	}
	return &vault, nil
}

func GetVaultsByUserId(pool *pgxpool.Pool, userId int64) ([]models.StoreVault, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT vaults.id, stores.name, vaults.name
		FROM vaults
		JOIN vaultaccess on vaultaccess.vault_id = vaults.id
		JOIN stores ON vaults.store_id = stores.id
		WHERE vaultaccess.user_id = $1
	`
	var rows, err = pool.Query(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var vaults []models.StoreVault = []models.StoreVault{}

	for rows.Next() {
		var v models.StoreVault

		err = rows.Scan(
			&v.ID,
			&v.StoreName,
			&v.Name,
		)

		if err != nil {
			return nil, err
		}

		vaults = append(vaults, v)
	}
	return vaults, nil
}

func GetVaultsByStoreUserId(pool *pgxpool.Pool, userId, storeId int64) ([]models.StoreVault, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT vaults.id, stores.name, vaults.name
		FROM vaults
		JOIN vaultaccess on vaultaccess.vault_id = vaults.id
		JOIN stores ON vaults.store_id = stores.id
		WHERE vaultaccess.user_id = $1 AND stores.id = $2
	`
	var rows, err = pool.Query(ctx, query, userId, storeId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var vaults []models.StoreVault = []models.StoreVault{}

	for rows.Next() {
		var v models.StoreVault

		err = rows.Scan(
			&v.ID,
			&v.StoreName,
			&v.Name,
		)

		if err != nil {
			return nil, err
		}

		vaults = append(vaults, v)
	}
	return vaults, nil
}

func UserHasVaultAcces(pool *pgxpool.Pool, userId, vaultId int64) (bool, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT EXISTS(
			SELECT 1
			FROM vaultaccess
			WHERE user_id = $1
			  AND vault_id = $2
		)`

	var exist bool

	err := pool.QueryRow(ctx, query, userId, vaultId).Scan(&exist)

	if err != nil {
		return false, err
	}

	return exist, nil
}

func GetAllVaults(pool *pgxpool.Pool) ([]models.StoreIdVault, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT vaults.id, stores.id, vaults.name
		FROM vaults
		JOIN stores ON vaults.store_id = stores.id
	`
	var rows, err = pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var vaults []models.StoreIdVault = []models.StoreIdVault{}

	for rows.Next() {
		var v models.StoreIdVault

		err = rows.Scan(
			&v.ID,
			&v.StoreId,
			&v.Name,
		)

		if err != nil {
			return nil, err
		}

		vaults = append(vaults, v)
	}
	return vaults, nil
}
