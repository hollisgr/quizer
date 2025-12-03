package user

import (
	"context"
	"log"
	"quizer_server/internal/db"
	"quizer_server/internal/model"
)

type Service interface {
	UserByLogin(ctx context.Context, login string) (model.User, error)
}

type userService struct {
	storage db.Storage
}

func New(s db.Storage) Service {
	return &userService{
		storage: s,
	}
}

func (s *userService) UserByLogin(ctx context.Context, login string) (model.User, error) {
	user, err := s.storage.UserByLogin(ctx, login)
	if err != nil {
		log.Println(err)
		return user, err
	}
	return user, nil
}
