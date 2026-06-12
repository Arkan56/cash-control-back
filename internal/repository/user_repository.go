package repository

import (
	"cash-control/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateWorkerUser(pool *pgxpool.Pool, req *models.CreateUserRequest) (*models.CreatedUserResponce, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `INSERT INTO users(
	userName,
	password,
	name,
	rol_id)
	values($1,$2,$3,$4)
	returning id, userName, name`

	var createdUser models.CreatedUserResponce

	err := pool.QueryRow(ctx, query, req.UserName, req.Password, req.Name, 2).Scan(
		&createdUser.ID,
		&createdUser.UserName,
		&createdUser.Name,
	)

	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}
