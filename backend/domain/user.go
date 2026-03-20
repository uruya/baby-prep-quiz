package domain

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

type UserRepository interface {
	Create(name, email, passwordHash string) (*User, error)
	FindByEmail(email string) (*User, string, error) // User, passwordHash, error
}
