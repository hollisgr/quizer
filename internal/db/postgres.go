package db

import (
	"context"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	UserByLogin(ctx context.Context, login string) (model.User, error)

	CreateGame(ctx context.Context, data dto.CreateNewGame) (int, error)
	GameList(ctx context.Context) ([]model.Game, error)
	GameLoad(ctx context.Context, id int) (model.Game, error)
	UpdateGame(ctx context.Context, updated model.Game) (int, error)

	CreateQuestion(ctx context.Context, data dto.CreateNewQuestionRequest) (int, error)
	QuestionLoad(ctx context.Context, id int) (model.Question, error)
	QuestionsByGameId(ctx context.Context, gameId int) ([]model.Question, error)
	UpdateQuestion(ctx context.Context, updated model.Question) (int, error)
	DeleteQuestion(ctx context.Context, id int) (int, error)
}

type storage struct {
	db *pgxpool.Pool
}

func New(p *pgxpool.Pool) Storage {
	return &storage{
		db: p,
	}
}
