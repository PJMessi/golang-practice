package model

type UserRes struct {
	Id        string  `json:"id"`
	Email     string  `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	CreatedAt string  `json:"createdAt"`
}
