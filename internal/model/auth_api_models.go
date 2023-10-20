package model

type LoginApiReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginApiRes struct {
	User UserRes `json:"user"`
	Jwt  string  `json:"jwt"`
}
