package models

// User model.
type User struct {
	// Identifier.
	ID int64 `db:"id"`
	// User name.
	Username string `db:"username"`
	// Password hash.
	PasswordHash string `db:"password_hash"`
}
