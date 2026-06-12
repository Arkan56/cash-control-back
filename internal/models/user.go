package models

type User struct {
	ID       int64  `db:"id"`
	UserName string `db:"userName"`
	Password string `db:"password"`
	Name     string `db:"name"`
	IdRol    int32  `db:"rol_id"`
}

type CreatedUserResponce struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	Name     string `json:"name"`
}

type CreateUserRequest struct {
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
