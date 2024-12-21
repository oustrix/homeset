package dto

import (
	"github.com/guregu/null/zero"
)

// GetUserInput is a set of filters to find a user.
type GetUserInput struct {
	IDEq       zero.Int
	UsernameEq zero.String
}

// CreateUserInput is a set of data for create new user.
type CreateUserInput struct {
	Username     string
	PasswordHash string
}
