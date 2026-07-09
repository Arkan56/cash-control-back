package repository

import (
	"cash-control/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateMovement(pool *pgxpool.Pool, detail string, amount int64, amountCategoryID int32, vaultID, userID int64) (*models.Movement, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var storeQuery string = `SELECT store_id FROM vaults WHERE id = $1`

	var storeId int64

	err = tx.QueryRow(ctx, storeQuery,
		vaultID).Scan(&storeId)

	if err != nil {
		return nil, err
	}

	var vaultLock string = `SELECT balance FROM vaults WHERE id = $1 FOR UPDATE `

	var vaulBalace int64
	err = tx.QueryRow(ctx, vaultLock, vaultID).Scan(&vaulBalace)

	if err != nil {
		return nil, err
	}

	if amount < 0 && amount*(-1) > vaulBalace {
		return nil, fmt.Errorf("cannot create movement: requested amount %d exceeds vault balance %d", amount, vaulBalace)
	}

	var query string = `
	INSERT INTO movements (
	detail, 
	amount, 
	amount_category_id, 
	vault_id, 
	user_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, detail, amount, amount_category_id, vault_id, user_id`

	var movement models.Movement

	err = tx.QueryRow(ctx, query,
		detail,
		amount,
		amountCategoryID,
		vaultID,
		userID).Scan(
		&movement.ID,
		&movement.CreatedAt,
		&movement.Detail,
		&movement.Amount,
		&movement.AmountCategoryID,
		&movement.VaultID,
		&movement.UserID,
	)

	if err != nil {
		return nil, err
	}

	var updateVault string = `UPDATE vaults SET balance = balance + $1 WHERE id = $2`

	_, err = tx.Exec(ctx, updateVault, amount, vaultID)

	if err != nil {
		return nil, err
	}

	var updateStore string = `UPDATE stores SET balance = balance + $1 WHERE id = $2`

	_, err = tx.Exec(ctx, updateStore, amount, storeId)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &movement, nil
}

func GetAllMovements(pool *pgxpool.Pool, vaultId int64) ([]models.MovementWithUser, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT movements.id, movements.created_at, movements.detail, movements.amount, users.name
	FROM movements 
	JOIN users ON movements.user_id = users.id
	WHERE movements.vault_id = $1
	`
	var rows, err = pool.Query(ctx, query, vaultId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movs []models.MovementWithUser = []models.MovementWithUser{}

	for rows.Next() {
		var m models.MovementWithUser

		err = rows.Scan(
			&m.ID,
			&m.CreatedAt,
			&m.Detail,
			&m.Amount,
			&m.UserName,
		)

		if err != nil {
			return nil, err
		}

		movs = append(movs, m)
	}

	return movs, nil
}
