package services

import (
	"quizer_server/internal/middleware"
	"quizer_server/internal/service/game"
	"quizer_server/internal/service/jwt"
	"quizer_server/internal/service/lobby"
	"quizer_server/internal/service/question"
	"quizer_server/internal/service/user"
)

type Services struct {
	UserSvc     user.Service
	GameSvc     game.Service
	LobbySvc    lobby.Service
	QuestionSvc question.Service
	JwtSvc      jwt.Service
	UserAuth    middleware.UserAuthenticator
}
