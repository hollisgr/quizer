package handler

import (
	"fmt"
	"net/http"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h *handler) QuestionById(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		sendError(c, http.StatusBadRequest, "question id is required")
		return
	}

	res, err := h.questionSvc.Load(c.Request.Context(), id)

	if err != nil {
		if err == pgx.ErrNoRows {
			sendError(c, http.StatusNotFound, "question not found")
			return
		}
		sendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	sendSuccess(c, http.StatusOK, res)
}

func (h *handler) QuestionsByGameId(c *gin.Context) {
	gameIdStr := c.Params.ByName("game_id")
	gameId := 0
	_, err := fmt.Sscanf(gameIdStr, "%d", &gameId)
	if err != nil {
		sendError(c, http.StatusBadRequest, "game_id is required")
		return
	}

	res, err := h.questionSvc.ListByGameId(c.Request.Context(), gameId)

	if err != nil {
		sendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	sendSuccess(c, http.StatusOK, res)
}

func (h *handler) CreateQuestion(c *gin.Context) {
	req := dto.CreateNewQuestionRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		sendError(c, http.StatusBadRequest, "body req err")
		return
	}

	id, err := h.questionSvc.Create(c.Request.Context(), req)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	resp := map[string]any{
		"id": id,
	}

	sendSuccess(c, http.StatusOK, resp)
}

func (h *handler) UpdateQuestion(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)

	if err != nil {
		sendError(c, http.StatusBadRequest, "question id is required")
		return
	}

	req := model.Question{}
	err = c.BindJSON(&req)
	if err != nil {
		sendError(c, http.StatusBadRequest, "body req err")
		return
	}

	id, err = h.questionSvc.Update(c.Request.Context(), req)
	if err != nil || id == 0 {
		sendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	resp := map[string]any{
		"id": id,
	}

	sendSuccess(c, http.StatusOK, resp)
}

func (h *handler) DeleteQuestion(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id := 0
	_, err := fmt.Sscanf(idStr, "%d", &id)

	if err != nil {
		sendError(c, http.StatusBadRequest, "question id is required")
		return
	}

	id, err = h.questionSvc.DeleteById(c.Request.Context(), id)
	if err != nil || id == 0 {
		sendError(c, http.StatusInternalServerError, "internal err")
		return
	}

	resp := map[string]any{
		"id": id,
	}

	sendSuccess(c, http.StatusOK, resp)
}
