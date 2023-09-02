package model

type UserRegApiReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegApiRes struct {
	User UserRes `json:"user"`
}
