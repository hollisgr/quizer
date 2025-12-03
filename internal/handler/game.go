package handler

import (
	"fmt"
	"net/http"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h *handler) CreateGame(c *gin.Context) {
	req := dto.CreateNewGameRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		sendError(c, http.StatusBadRequest, "body req err")
		return
	}

	data := dto.CreateNewGame{
		OwnerId:     h.jwtSvc.IDFromToken(c.Value("access_token").(string)),
		Description: req.Description,
		Link:        req.Link,
	}

	id, err := h.gameSvc.CreateNewGame(c.Request.Context(), data)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	resp := map[string]any{
		"id": id,
	}

	sendSuccess(c, http.StatusOK, resp)
}

func (h *handler) GameList(c *gin.Context) {
	list, err := h.gameSvc.GameList(c.Request.Context())

	if err != nil {
		if err == pgx.ErrNoRows {
			sendError(c, http.StatusNotFound, "game list is empty")
			return
		}
		sendError(c, http.StatusInternalServerError, "internal err")
		return

	}

	sendSuccess(c, http.StatusOK, list)
}

func (h *handler) GameLoad(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)

	if err != nil || id == 0 {
		sendError(c, http.StatusBadRequest, "incorrect game_id")
		return
	}

	res, err := h.gameSvc.GameLoad(c.Request.Context(), id)

	if err != nil {
		if err == pgx.ErrNoRows {
			sendError(c, http.StatusNotFound, "game not found")
			return
		}
		sendError(c, http.StatusInternalServerError, "internal err")
		return

	}

	sendSuccess(c, http.StatusOK, res)
}

func (h *handler) UpdateGame(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)
	req := model.Game{
		Id: id,
	}
	err = c.BindJSON(&req)
	if err != nil {
		sendError(c, http.StatusBadRequest, "body req err")
		return
	}

	id, err = h.gameSvc.UpdateGame(c.Request.Context(), req)
	if err != nil || id == 0 {
		sendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	resp := map[string]any{
		"id": id,
	}

	sendSuccess(c, http.StatusOK, resp)
}

func (h *handler) DeleteGame(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)

	if err != nil {
		sendError(c, http.StatusBadRequest, "game id is required")
		return
	}

	id, err = h.gameSvc.DeleteGame(c.Request.Context(), id)
	if err != nil || id == 0 {
		sendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	resp := map[string]any{
		"id": id,
	}

	sendSuccess(c, http.StatusOK, resp)
}
