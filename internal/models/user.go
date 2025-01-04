package models

// User model.
type User struct {
	// User name.
	Username string `db:"username"`
	// Password hash.
	PasswordHash string `db:"password_hash"`
}
