package database

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type UserRepository interface {
	Create(email string) (User, error)
}
