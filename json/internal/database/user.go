package database

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserParams struct {
	Email    string
	Password string
}

type UpdateUserParams struct {
	Email    string
	Password string
}

type UserRepository interface {
	Create(params CreateUserParams) (User, error)
	GetByEmail(email string) (User, error)
	GetByID(id int) (User, error)
	Update(id int, params UpdateUserParams) (User, error)
}
