package domain

import (
	"context"

	"github.com/google/uuid"
)

type Player struct {
	UUID      uuid.UUID `json:"uuid" db:"uuid"`
	UserName  string    `json:"user_name" db:"user_name"`
	LobbyUUID uuid.UUID `json:"lobby_uuid" db:"lobby_uuid"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	GameId    int       `json:"game_id" db:"game_id"`
}

type PlayerRepository interface {
	PlayersByGameUUID(ctx context.Context, gameUUID uuid.UUID) ([]Player, error)
	SavePlayer(ctx context.Context, newPlayer Player) error
}
