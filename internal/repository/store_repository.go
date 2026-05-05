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
