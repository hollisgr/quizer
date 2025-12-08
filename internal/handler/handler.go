package handler

import (
	"fmt"
	"net/http"
	"quizer_server/internal/app/services"
	"quizer_server/internal/config"
	"quizer_server/internal/middleware"
	"quizer_server/internal/service/game"
	"quizer_server/internal/service/jwt"
	"quizer_server/internal/service/question"
	"quizer_server/internal/service/user"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler interface {
	Register()
}

type handler struct {
	router            *gin.Engine
	userSvc           user.Service
	gameSvc           game.Service
	questionSvc       question.Service
	jwtSvc            jwt.Service
	userAuth          middleware.UserAuthenticator
	updater           websocket.Upgrader
	activeConnections map[uuid.UUID]*websocket.Conn
}

func New(r *gin.Engine, s services.Services) Handler {
	return &handler{
		router:      r,
		userSvc:     s.UserSvc,
		jwtSvc:      s.JwtSvc,
		userAuth:    s.UserAuth,
		gameSvc:     s.GameSvc,
		questionSvc: s.QuestionSvc,
		updater: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Register configures HTTP routes for managing wallet resources.
func (h *handler) Register() {

	cfg := config.GetConfig()
	fmt.Println(cfg.CORS.AllowedOrigins)

	configCORS := cors.DefaultConfig()
	configCORS.AllowOrigins = cfg.CORS.AllowedOrigins
	configCORS.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	configCORS.AllowCredentials = true

	h.router.Use(cors.New(configCORS))

	protected := h.router.Group("/", h.userAuth.Authorization())

	h.router.GET("/login", h.Login)
	h.router.GET("/ws/:player_uuid", h.wsHandler)

	protected.GET("/user/:login", h.UserByLogin)

	protected.GET("/questions/:id", h.QuestionById)
	protected.GET("/questions/game/:game_id", h.QuestionsByGameId)
	protected.POST("/questions", h.CreateQuestion)
	protected.POST("/questions/:id", h.UpdateQuestion)
	protected.DELETE("/questions/:id", h.DeleteQuestion)

	protected.GET("/games", h.GameList)
	protected.GET("/games/:id", h.GameLoad)
	protected.POST("/games/:id", h.UpdateGame)
	protected.POST("/games", h.CreateGame)
	protected.DELETE("/games/:id", h.DeleteGame)
}

// sendError sends an error response to the client with a specified HTTP status code and error message.
func sendError(c *gin.Context, code int, message any) {
	c.AbortWithStatusJSON(code, gin.H{
		"success": false,
		"message": fmt.Sprint(message),
	})
}

// sendSuccess sends a success response to the client with a specified HTTP status code and success message/data.
func sendSuccess(c *gin.Context, code int, message any) {
	c.JSON(code, gin.H{
		"success": true,
		"message": message,
	})
}

func parseAuthHeader(header string) (string, error) {
	if header == "" {
		return header, fmt.Errorf("access denied, Authorization header required")
	}
	headerArr := strings.Split(header, " ")
	if headerArr[0] != "Basic" {
		return header, fmt.Errorf("access denied, Basic64 authorization required")
	}
	return headerArr[1], nil
}
