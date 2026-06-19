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
	UserName string `json:"user_name" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}
