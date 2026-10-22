package repository

import (
	"cash-control/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateAmountCategory(pool *pgxpool.Pool, description string) (*models.AmountCategory, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `INSERT INTO amountCategory (description)
	VALUES ($1)
	RETURNING id, description
	`

	var amountCategory models.AmountCategory

	err := pool.QueryRow(ctx, query, description).Scan(&amountCategory.ID, &amountCategory.Description)

	if err != nil {
		return nil, err
	}

	return &amountCategory, nil
}

func GetAllamountCategories(pool *pgxpool.Pool) ([]models.AmountCategory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query = `SELECT * FROM amountCategory`

	var rows, err = pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var amountCategories []models.AmountCategory = []models.AmountCategory{}

	for rows.Next() {
		var amountCategory models.AmountCategory

		err = rows.Scan(
			&amountCategory.ID,
			&amountCategory.Description,
		)

		if err != nil {
			return nil, err
		}

		amountCategories = append(amountCategories, amountCategory)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return amountCategories, err
}
