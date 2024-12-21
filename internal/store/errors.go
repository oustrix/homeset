package store

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("entity not found")
	ErrUniqueViolation = errors.New("unique violation")
)
