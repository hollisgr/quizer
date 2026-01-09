package apiV1

import (
	"quizer_server/internal/usecases"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameUC usecases.GameUseCases
}

func (h *GameHandler) CreateGame(c *gin.Context) {
	// req := dto.CreateNewGameRequest{}
	// err := c.BindJSON(&req)
	// if err != nil {
	// 	responce.SendError(c, http.StatusBadRequest, "body req err")
	// 	return
	// }

	// data := dto.CreateNewGame{
	// 	OwnerId:     h.jwtSvc.IDFromToken(c.Value("access_token").(string)),
	// 	Description: req.Description,
	// 	Link:        req.Link,
	// }

	// id, err := h.gameSvc.CreateNewGame(c.Request.Context(), data)
	// if err != nil {
	// 	responce.SendError(c, http.StatusInternalServerError, "internal err")
	// 	return
	// }

	// resp := map[string]any{
	// 	"id": id,
	// }

	// responce.SendSuccess(c, http.StatusOK, resp)
}

func (h *GameHandler) GameList(c *gin.Context) {
	// list, err := h.gameSvc.GameList(c.Request.Context())

	// if err != nil {
	// 	if err == pgx.ErrNoRows {
	// 		responce.SendError(c, http.StatusNotFound, "game list is empty")
	// 		return
	// 	}
	// 	responce.SendError(c, http.StatusInternalServerError, "internal err")
	// 	return

	// }

	// responce.SendSuccess(c, http.StatusOK, list)
}

func (h *GameHandler) GameLoad(c *gin.Context) {
	// idStr := c.Params.ByName("id")
	// id := 0
	// _, err := fmt.Sscanf(idStr, "%d", &id)

	// if err != nil || id == 0 {
	// 	responce.SendError(c, http.StatusBadRequest, "incorrect game_id")
	// 	return
	// }

	// res, err := h.gameSvc.GameLoad(c.Request.Context(), id)

	// if err != nil {
	// 	if err == pgx.ErrNoRows {
	// 		responce.SendError(c, http.StatusNotFound, "game not found")
	// 		return
	// 	}
	// 	responce.SendError(c, http.StatusInternalServerError, "internal err")
	// 	return

	// }

	// responce.SendSuccess(c, http.StatusOK, res)
}

func (h *GameHandler) UpdateGame(c *gin.Context) {
	// idStr := c.Params.ByName("id")
	// id := 0
	// _, err := fmt.Sscanf(idStr, "%d", &id)
	// req := model.Game{
	// 	Id: id,
	// }
	// err = c.BindJSON(&req)
	// if err != nil {
	// 	responce.SendError(c, http.StatusBadRequest, "body req err")
	// 	return
	// }

	// id, err = h.gameSvc.UpdateGame(c.Request.Context(), req)
	// if err != nil || id == 0 {
	// 	responce.SendError(c, http.StatusInternalServerError, "internal err")
	// 	return
	// }

	// resp := map[string]any{
	// 	"id": id,
	// }

	// responce.SendSuccess(c, http.StatusOK, resp)
}

func (h *GameHandler) DeleteGame(c *gin.Context) {
	// idStr := c.Params.ByName("id")
	// id := 0
	// _, err := fmt.Sscanf(idStr, "%d", &id)

	// if err != nil {
	// 	responce.SendError(c, http.StatusBadRequest, "game id is required")
	// 	return
	// }

	// id, err = h.gameSvc.DeleteGame(c.Request.Context(), id)
	// if err != nil || id == 0 {
	// 	responce.SendError(c, http.StatusInternalServerError, "internal err")
	// 	return
	// }

	// resp := map[string]any{
	// 	"id": id,
	// }

	// responce.SendSuccess(c, http.StatusOK, resp)
}
