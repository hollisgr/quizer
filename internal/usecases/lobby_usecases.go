package usecases

import "quizer_server/internal/domain"

type LobbyUseCases struct {
	lobbyRepo domain.LobbyRepository
}

func NewLobbyUseCases(lobbyRepo domain.LobbyRepository) *LobbyUseCases {
	return &LobbyUseCases{
		lobbyRepo: lobbyRepo,
	}
}
