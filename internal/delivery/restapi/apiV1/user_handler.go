package apiV1

import (
	"errors"
	"net/http"
	"quizer_server/internal/delivery/restapi/responce"
	"quizer_server/internal/domain"
	"quizer_server/internal/usecases"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC usecases.UserUseCases
}

func NewUserHandler(uuc usecases.UserUseCases) *UserHandler {
	return &UserHandler{
		userUC: uuc,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	login, password, ok := c.Request.BasicAuth()
	if !ok {
		responce.SendError(c, http.StatusUnauthorized, "basic auth requested")
		return
	}

	token, err := h.userUC.Login(c.Request.Context(), login, password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			responce.SendError(c, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidCredentials) {
			responce.SendError(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		responce.SendError(c, http.StatusInternalServerError, "internal error")
		return
	}
	responce.SendSuccess(c, http.StatusOK, token)
}
