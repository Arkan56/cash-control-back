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

func GetUserById(pool *pgxpool.Pool, id int64) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var query = `SELECT id, userName, name, password, rol_id FROM users WHERE users.id = $1`

	var user models.User
	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.UserName,
		&user.Name,
		&user.Password,
		&user.IdRol,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers(pool *pgxpool.Pool) ([]models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, userName, name, password, rol_id FROM users
	`

	var rows, err = pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User = []models.User{}

	for rows.Next() {
		var user models.User

		err = rows.Scan(
			&user.ID,
			&user.UserName,
			&user.Name,
			&user.Password,
			&user.IdRol,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, err

}

func GetUserByUserName(pool *pgxpool.Pool, user_name string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	var query = `SELECT id, userName, name, password, rol_id FROM users WHERE users.userName = $1`

	var user models.User
	err := pool.QueryRow(ctx, query, user_name).Scan(
		&user.ID,
		&user.UserName,
		&user.Name,
		&user.Password,
		&user.IdRol,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateVaultAccess(pool *pgxpool.Pool, userId, vaultId int64) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		INSERT INTO vaultAccess (vault_id, user_id)
		VALUES ($1, $2)
	`

	_, err := pool.Exec(ctx, query, vaultId, userId)
	return err
}

func CreateStoreAccess(pool *pgxpool.Pool, userId, storeId int64) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		INSERT INTO storeAccess (store_id, user_id)
		VALUES ($1, $2)
	`

	_, err := pool.Exec(ctx, query, storeId, userId)
	return err
}
func GetVaultAccessByUserId(pool *pgxpool.Pool, userId int64) ([]models.VaultAccessResponse, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT vault_id
		FROM vaultAccess
		WHERE user_id = $1
	`
	var rows, err = pool.Query(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var vaults []models.VaultAccessResponse = []models.VaultAccessResponse{}

	for rows.Next() {
		var v models.VaultAccessResponse

		err = rows.Scan(
			&v.VaultID,
		)

		if err != nil {
			return nil, err
		}

		vaults = append(vaults, v)
	}
	return vaults, nil
}

func GetStoresAccessByUserId(pool *pgxpool.Pool, userId int64) ([]models.StoreAccessResponse, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT store_id
		FROM storeAccess
		WHERE user_id = $1
	`
	var rows, err = pool.Query(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stores []models.StoreAccessResponse = []models.StoreAccessResponse{}

	for rows.Next() {
		var s models.StoreAccessResponse

		err = rows.Scan(
			&s.StoreID,
		)

		if err != nil {
			return nil, err
		}

		stores = append(stores, s)
	}
	return stores, nil
}

func SyncAccessByUserId(pool *pgxpool.Pool, userId int64, storeIDs, vaultIDs []int64) (bool, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	//stores
	if _, err := tx.Exec(ctx,
		`DELETE FROM storeaccess WHERE user_id = $1 AND store_id != ALL($2)`,
		userId, storeIDs,
	); err != nil {
		return false, err
	}

	for _, storeId := range storeIDs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO storeaccess (store_id, user_id)
				 VALUES ($1, $2)
				 ON CONFLICT (store_id, user_id) DO NOTHING`,
			storeId, userId,
		); err != nil {
			return false, err
		}
	}

	//vaults
	if _, err := tx.Exec(ctx,
		`DELETE FROM vaultaccess WHERE user_id = $1 AND vault_id != ALL($2)`,
		userId, vaultIDs,
	); err != nil {
		return false, err
	}

	for _, vaultId := range vaultIDs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO vaultaccess (vault_id, user_id)
				 VALUES ($1, $2)
				 ON CONFLICT (vault_id, user_id) DO NOTHING`,
			vaultId, userId,
		); err != nil {
			return false, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}
