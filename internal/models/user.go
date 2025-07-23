package models

import (
	"gitlab.plainsurf.com/plainsurf/poc/jumed/poc-project/internal/config"
)

type User struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	Role      string `json:"role" db:"role"`
	CreatedAt string `json:"created_at" db:"created_at"` // Handles nullable datetime correctly
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(u *User) error {
	_, err := config.DB.Exec(`
		INSERT INTO users (name, email, password, role)
		VALUES (?, ?, ?, ?)`,
		u.Name, u.Email, u.Password, u.Role,
	)
	return err
}

func GetUserByEmail(email string) (User, error) {
	var user User
	err := config.DB.Get(&user, `
		SELECT id, name, email, password, role, created_at
		FROM users WHERE email = ?`, email,
	)
	return user, err
}
