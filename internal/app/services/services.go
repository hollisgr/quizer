package services

import (
	"quizer_server/internal/infrastracture/jwt"
	"quizer_server/internal/usecases"
)

type UseCases struct {
	UserUseCase  usecases.UserUseCases
	TokenManager *jwt.Manager
	// GameSvc     game.Service
	// LobbySvc    lobby.Service
	// QuestionSvc question.Service
	// JwtSvc   jwt.Service
	// UserAuth middleware.UserAuthenticator
}
