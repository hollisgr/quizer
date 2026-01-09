package usecases

import (
	"context"
	"quizer_server/internal/domain"
)

type TokenGenerator interface {
	CreateToken(userID int, login string) (string, error)
}

type UserUseCases interface {
	Login(ctx context.Context, login string, password string) (string, error)
}

type userUseCases struct {
	userRepo       domain.UserRepository
	tokenGenerator TokenGenerator
}

func NewUserUseCases(repo domain.UserRepository, tg TokenGenerator) UserUseCases {
	return &userUseCases{
		userRepo:       repo,
		tokenGenerator: tg,
	}
}

func (uc *userUseCases) Login(ctx context.Context, login string, password string) (string, error) {
	user, err := uc.userRepo.ByLogin(ctx, login)
	if err != nil {
		return "", domain.ErrUserNotFound
	}

	if password != user.Password {
		return "", domain.ErrInvalidCredentials
	}

	return uc.tokenGenerator.CreateToken(user.Id, user.Login)
}
