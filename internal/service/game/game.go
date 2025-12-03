package game

import (
	"context"
	"log"
	"quizer_server/internal/db"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"
)

type Service interface {
	CreateNewGame(ctx context.Context, data dto.CreateNewGame) (int, error)
	GameList(ctx context.Context) ([]model.Game, error)
	GameLoad(ctx context.Context, id int) (model.Game, error)
	UpdateGame(ctx context.Context, updated model.Game) (int, error)
}

type gameService struct {
	storage db.Storage
}

func New(s db.Storage) Service {
	return &gameService{
		storage: s,
	}
}

func (gs *gameService) CreateNewGame(ctx context.Context, data dto.CreateNewGame) (int, error) {
	id, err := gs.storage.CreateGame(ctx, data)
	if err != nil {
		log.Println(err)
		return id, err
	}
	return id, err
}

func (gs *gameService) GameList(ctx context.Context) ([]model.Game, error) {
	list, err := gs.storage.GameList(ctx)
	if err != nil {
		log.Println(err)
		return list, err
	}
	return list, nil
}

func (gs *gameService) GameLoad(ctx context.Context, id int) (model.Game, error) {
	res, err := gs.storage.GameLoad(ctx, id)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, nil
}

func (gs *gameService) UpdateGame(ctx context.Context, updated model.Game) (int, error) {
	res, err := gs.storage.UpdateGame(ctx, updated)
	if err != nil || res == 0 {
		log.Println(err)
		return res, err
	}
	return res, nil
}
