package services

import (
	"quizer_server/internal/infrastracture/jwt"
	"quizer_server/internal/usecases"
)

type UseCases struct {
	UserUseCase *usecases.UserUseCases
	// LobbyUseCase    usecases.LobbyUseCases
	// GameUseCase     usecases.GameUseCases
	// QuestionUseCase usecases.QuestionUseCases
	TokenManager *jwt.Manager
}
