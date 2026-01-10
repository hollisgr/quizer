package usecases

import (
	"context"
	"log"
	"quizer_server/internal/domain"
)

type TokenGenerator interface {
	CreateToken(userID int, login string) (string, error)
}

type UserUseCases struct {
	userRepo       domain.UserRepository
	tokenGenerator TokenGenerator
}

func NewUserUseCases(repo domain.UserRepository, tg TokenGenerator) *UserUseCases {
	return &UserUseCases{
		userRepo:       repo,
		tokenGenerator: tg,
	}
}

func (uc *UserUseCases) Login(ctx context.Context, login string, password string) (string, error) {
	user, err := uc.userRepo.ByLogin(ctx, login)
	if err != nil {
		log.Println("user usecase login err:", err)
		return "", domain.ErrUserNotFound
	}

	if password != user.Password {
		log.Println("user usecase login err: incorrect creds")
		return "", domain.ErrInvalidCredentials
	}

	return uc.tokenGenerator.CreateToken(user.Id, user.Login)
}
