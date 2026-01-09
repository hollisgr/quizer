package apiV1

import (
	"log"
	"net/http"
	"quizer_server/internal/delivery/restapi/responce"
	"quizer_server/internal/domain"
	"quizer_server/internal/usecases"

	"github.com/gin-gonic/gin"
)

type LobbyHandler struct {
	lobbyUC usecases.LobbyUseCases
}

type LobbyUseCases interface {
	// описание юзкейсов, которые нужны для работы лобби хэндлера
}

func (h *LobbyHandler) CreateLobby(c *gin.Context) {
	req := domain.Lobby{}
	err := c.BindJSON(&req)
	if err != nil {
		responce.SendError(c, http.StatusBadRequest, "body req err")
		log.Println("create lobby bind json err:", err)
		return
	}
	// count, err := h.lobbySvc.Create(c.Request.Context(), req)
	// if err != nil {
	// 	log.Println("handler create new lobby err:", err)
	// }
	// responce.SendSuccess(c, http.StatusOK, gin.H{
	// 	"success":         true,
	// 	"questions_count": count,
	// })
}

func (h *LobbyHandler) LobbyList(c *gin.Context) {
	// res, err := h.lobbySvc.List(c.Request.Context())
	// if err != nil {
	// 	log.Println("handler lobby list err:", err)
	// }
	// responce.SendSuccess(c, http.StatusOK, res)
}
