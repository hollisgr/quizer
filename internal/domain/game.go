package domain

import (
	"context"
	"time"
)

type Game struct {
	Id          int       `json:"game_id" db:"id"`
	Description string    `json:"description" db:"description"`
	Owner       int       `json:"owner_id" db:"owner_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Link        string    `json:"link" db:"link"`
}

type GameWithLogin struct {
	Game
	Login string `json:"login" db:"login"`
}

type GameRepository interface {
	CreateGame(ctx context.Context, data Game) (int, error)
	GameList(ctx context.Context) ([]GameWithLogin, error)
	GameLoad(ctx context.Context, id int) (GameWithLogin, error)
	UpdateGame(ctx context.Context, updated Game) (int, error)
	UpdateFilePath(ctx context.Context, gameId int, path string) (int, error)
	DeleteGame(ctx context.Context, id int) (int, error)
}
