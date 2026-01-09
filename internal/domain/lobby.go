package domain

import (
	"context"

	"github.com/google/uuid"
)

type Lobby struct {
	UUID      uuid.UUID `json:"uuid" db:"uuid"`
	GameId    int       `json:"game_id" db:"game_id"`
	IsStarted bool      `json:"is_started" db:"is_started"`
}

type LobbyRepository interface {
	CreateLobby(ctx context.Context, data Lobby) error
	LobbyLoadByUUID(ctx context.Context, uuid uuid.UUID) (Lobby, error)
	UpdateLobby(ctx context.Context, lobbyUUID uuid.UUID) error
	LobbyList(ctx context.Context) ([]Lobby, error)
}
