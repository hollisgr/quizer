package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"quizer_server/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *handler) Login(c *gin.Context) {
	data, err := parseAuthHeader(c.GetHeader("Authorization"))
	if err != nil {
		sendError(c, http.StatusUnauthorized, fmt.Sprint(err))
		return
	}

	body, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		sendError(c, http.StatusUnauthorized, "Access denied, base64 decoding error")
		return
	}

	dataBytes := bytes.Split(body, []byte(":"))
	login := string(dataBytes[0])
	password := string(dataBytes[1])

	user, err := h.userSvc.UserByLogin(c, login)

	if err != nil {
		sendError(c, http.StatusUnauthorized, "Access denied, user not found")
		return
	}

	if user.Password != password {
		sendError(c, http.StatusUnauthorized, "Access denied, password do not match")
		return
	}

	tokens := h.jwtSvc.CreateToken(c.Request.Context(), model.JwtRequest{
		Login:    login,
		Password: password,
	})

	sendSuccess(c, http.StatusOK, tokens)
}

func (h *handler) UserByLogin(c *gin.Context) {
	login := c.Param("login")
	user, err := h.userSvc.UserByLogin(c, login)

	if err != nil {
		sendError(c, http.StatusNotFound, "user not found")
		return
	}

	sendSuccess(c, http.StatusOK, user)
}
