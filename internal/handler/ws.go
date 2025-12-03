package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// wsHandler manages WebSocket connections for real-time communication between users.
// It upgrades the HTTP connection to WebSocket, registers the user connection, reads messages,
// and broadcasts responses to both players involved in a game session.
func (h *handler) wsHandler(c *gin.Context) {
	paramUUID := c.Param("player_uuid")
	if paramUUID == "" {
		sendError(c, http.StatusUnauthorized, "player uuid is empty")
		return
	}

	playerUUID, err := uuid.Parse(paramUUID)
	if err != nil {
		sendError(c, http.StatusUnauthorized, "player uuid is incorrect")
		return
	}

	ws, err := h.updater.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "ws error")
		return
	}

	log.Println("player connected: ", playerUUID)

	defer func() {
		ws.Close()
	}()

	h.wsRegistration(playerUUID, ws)

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		parseMsg(msg, msgType)
		h.activeConnections[playerUUID].WriteJSON(gin.H{"conn": true})
	}
}

// wsRegistration adds a new WebSocket connection to the active connections map indexed by user UUID.
func (h *handler) wsRegistration(playerUUID uuid.UUID, conn *websocket.Conn) {
	h.activeConnections[playerUUID] = conn
}

func parseMsg(msg []byte, msgType int) {}
