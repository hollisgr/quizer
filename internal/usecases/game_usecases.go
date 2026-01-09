package usecases

import "quizer_server/internal/domain"

type GameUseCases interface{}

type gameUseCases struct {
	gameRepo domain.GameRepository
}

func NewGameUseCases(gameRepo domain.GameRepository) GameUseCases {
	return &gameUseCases{
		gameRepo: gameRepo,
	}
}

func (uc *gameUseCases) Test() {}
