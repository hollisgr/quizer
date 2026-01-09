package domain

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	Id       int    `json:"user_id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

type UserRepository interface {
	ByLogin(ctx context.Context, login string) (User, error)
}
