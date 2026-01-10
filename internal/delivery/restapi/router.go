package restapi

import (
	"quizer_server/internal/app/services"
	"quizer_server/internal/delivery/restapi/apiV1"
	"quizer_server/internal/delivery/restapi/middleware"

	"github.com/gin-gonic/gin"
)

type router struct {
	router       *gin.Engine
	userHandler  *apiV1.UserHandler
	gameHandler  *apiV1.GameHandler
	tokenManager middleware.TokenManager
}

func NewRouter(uc services.UseCases, r *gin.Engine) *router {
	return &router{
		router:       r,
		userHandler:  apiV1.NewUserHandler(uc.UserUseCase),
		tokenManager: uc.TokenManager,
	}
}

func (r *router) Register() {
	publicV1 := r.router.Group("/api/v1/public")
	publicV1.GET("/login", r.userHandler.Login)

	protectedV1 := r.router.Group("/api/v1/protected")
	protectedV1.Use(middleware.AuthMiddleware(r.tokenManager))
	protectedV1.GET("/test")
}
