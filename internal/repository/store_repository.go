package repository

import (
	"cash-control/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateStore(pool *pgxpool.Pool, name string) (*models.Store, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
			INSERT INTO stores (name)
			VALUES ($1)
			RETURNING id, name, balance
	`

	var store models.Store

	var err error = pool.QueryRow(ctx, query, name).Scan(
		&store.ID,
		&store.Name,
		&store.Balance,
	)

	if err != nil {
		return nil, err
	}

	return &store, nil
}

func GetAllStores(pool *pgxpool.Pool) ([]models.Store, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, name, balance
		FROM stores
	`

	var rows, err = pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stores []models.Store = []models.Store{}

	for rows.Next() {
		var store models.Store

		err = rows.Scan(
			&store.ID,
			&store.Name,
			&store.Balance,
		)

		if err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stores, err

}

func GetStoresByUserId(pool *pgxpool.Pool, userId int64) ([]models.Store, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT stores.id, stores.name, stores.balance
		FROM stores
		JOIN storeaccess on storeaccess.store_id = stores.id
		WHERE storeaccess.user_id = $1
	`
	var rows, err = pool.Query(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stores []models.Store = []models.Store{}

	for rows.Next() {
		var s models.Store

		err = rows.Scan(
			&s.ID,
			&s.Name,
			&s.Balance,
		)

		if err != nil {
			return nil, err
		}

		stores = append(stores, s)
	}
	return stores, nil
}

func UserHasStoreAcces(pool *pgxpool.Pool, userId, storeId int64) (bool, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT EXISTS(
			SELECT 1
			FROM storeaccess
			WHERE user_id = $1
			  AND store_id = $2
		)`

	var exist bool

	err := pool.QueryRow(ctx, query, userId, storeId).Scan(&exist)

	if err != nil {
		return false, err
	}

	return exist, nil
}
