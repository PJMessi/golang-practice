package responses

import "time"

type UserResponse struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
}
