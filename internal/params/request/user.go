package request
// User represents a user in the system
type User struct {
	Username  string    `json:"username" validate:"required,min=3,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
}
type UserUpdate struct {
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
}
