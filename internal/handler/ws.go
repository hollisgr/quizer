package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"quizer_server/internal/model"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// wsHandler manages WebSocket connections for real-time communication between users.
// It upgrades the HTTP connection to WebSocket, registers the user connection, reads messages,
// and broadcasts responses to both players involved in a game session.
func (h *handler) wsHandler(c *gin.Context) {
	paramPlayerUUID := c.Query("player_uuid")
	paramLobbyUUID := c.Query("lobby_uuid")
	paramPlayerName := c.Query("player_name")
	paramIsAuth := c.Query("is_auth")
	isAdmin := false

	if paramPlayerUUID == "" {
		sendError(c, http.StatusBadRequest, "player uuid required")
		log.Println("player uuid is required")
		return
	}

	if paramLobbyUUID == "" {
		sendError(c, http.StatusBadRequest, "lobby uuid is required")
		log.Println("lobby uuid is required")
		return
	}

	if paramPlayerName == "" {
		sendError(c, http.StatusBadRequest, "player user name is required")
		log.Println("player user name is required")
		return
	}

	if paramIsAuth == "true" {
		isAdmin = true
	}

	playerUUID, err := uuid.Parse(paramPlayerUUID)
	if err != nil {
		sendError(c, http.StatusBadRequest, "player uuid is incorrect")
		log.Println("player uuid is incorrect:", paramPlayerUUID)
		return
	}

	lobbyUUID, err := uuid.Parse(paramLobbyUUID)
	if err != nil {
		sendError(c, http.StatusBadRequest, "game uuid is incorrect")
		log.Println("lobby uuid is incorrect:", paramLobbyUUID)
		return
	}

	ws, err := h.updater.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "ws error")
		return
	}

	log.Println("player connected:", playerUUID, "is auth:", isAdmin)

	defer func() {
		delete(h.sessions.activeConnections[lobbyUUID], playerUUID)
		h.updateUserList(lobbyUUID)
		log.Println("player disconnected:", playerUUID)
		ws.Close()
	}()

	data := PlayerData{
		Connection: ws,
		UserName:   paramPlayerName,
		IsAdmin:    isAdmin,
	}

	h.wsRegistration(c.Request.Context(), lobbyUUID, playerUUID, data)
	h.updateUserList(lobbyUUID)

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		h.parseMsg(c.Request.Context(), playerUUID, lobbyUUID, msg, msgType)
	}
}

// wsRegistration adds a new WebSocket connection to the active connections map indexed by user UUID.
func (h *handler) wsRegistration(ctx context.Context, lobbyUUID uuid.UUID, playerUUID uuid.UUID, data PlayerData) {
	_, ok := h.sessions.activeConnections[lobbyUUID]
	if !ok {
		h.sessions.activeConnections[lobbyUUID] = make(map[uuid.UUID]PlayerData)
	}

	h.sessions.mu.Lock()
	h.sessions.activeConnections[lobbyUUID][playerUUID] = data
	h.sessions.mu.Unlock()
	newPlayer := model.Player{
		UUID:      playerUUID,
		LobbyUUID: lobbyUUID,
		IsAdmin:   data.IsAdmin,
		UserName:  data.UserName,
	}
	h.gameSvc.SavePlayer(ctx, newPlayer)
}

func (h *handler) updateUserList(lobbyUUID uuid.UUID) {
	log.Println("update userlist")
	players := ""
	h.sessions.mu.Lock()
	for _, l := range h.sessions.activeConnections[lobbyUUID] {
		players += l.UserName + "/"
	}
	for _, l := range h.sessions.activeConnections[lobbyUUID] {
		l.Connection.WriteJSON(gin.H{
			"type": "lobby",
			"data": players,
		})
	}
	h.sessions.mu.Unlock()
}

func (h *handler) isAdmin(playerUUID, lobbyUUID uuid.UUID) bool {
	return h.sessions.activeConnections[lobbyUUID][playerUUID].IsAdmin
}

func (h *handler) parseMsg(ctx context.Context, playerUUID, lobbyUUID uuid.UUID, msg []byte, msgType int) {
	log.Println("lobby_uuid:", lobbyUUID, "player_uuid:", playerUUID, "msgType:", msgType, "msg:", string(msg))

	if string(msg) == "start" {
		h.sessions.mu.Lock()
		if h.isAdmin(playerUUID, lobbyUUID) {
			h.sessions.activeConnections[lobbyUUID][lobbyUUID].Connection.WriteJSON(gin.H{
				"type": "start",
			})
		} else {
			h.sessions.activeConnections[lobbyUUID][lobbyUUID].Connection.WriteJSON(gin.H{
				"type": "access_denied",
			})
		}
		h.sessions.mu.Unlock()
	}

	if string(msg) == "next" {
		h.sessions.mu.Lock()
		if h.isAdmin(playerUUID, lobbyUUID) {
			h.sessions.activeConnections[lobbyUUID][lobbyUUID].Connection.WriteJSON(gin.H{
				"type": "next",
			})
		} else {
			h.sessions.activeConnections[lobbyUUID][lobbyUUID].Connection.WriteJSON(gin.H{
				"type": "access_denied",
			})
		}
		h.sessions.mu.Unlock()
	}

	if string(msg) == "finish_quiz" {
		h.sessions.mu.Lock()
		if h.isAdmin(playerUUID, lobbyUUID) {
			h.sessions.activeConnections[lobbyUUID][lobbyUUID].Connection.WriteJSON(gin.H{
				"type": "finish_quiz",
			})
		} else {
			h.sessions.activeConnections[lobbyUUID][lobbyUUID].Connection.WriteJSON(gin.H{
				"type": "access_denied",
			})
		}
		h.sessions.mu.Unlock()
	}

	// if string(msg) == "final_question" {
	// 	h.sessions.mu.Lock()
	// 	for _, l := range h.sessions.activeConnections[lobbyUUID] {
	// 		l.Connection.WriteJSON(gin.H{
	// 			"type": "final_question",
	// 		})
	// 	}
	// 	h.sessions.mu.Unlock()
	// 	return
	// }

	// if string(msg) == "questions" {
	// 	lobby, _ := h.lobbySvc.LoadByUUID(ctx, lobbyUUID)
	// 	questions, _ := h.questionSvc.ListByGameId(ctx, lobby.GameId)
	// 	h.sessions.mu.Lock()
	// 	h.sessions.activeConnections[lobbyUUID][playerUUID].Connection.WriteJSON(gin.H{
	// 		"type": "questions",
	// 		"data": questions})
	// 	h.sessions.mu.Unlock()
	// 	return
	// }

	if strings.Contains(string(msg), "start_lobby") {
		lobby, _ := h.lobbySvc.LoadByUUID(ctx, lobbyUUID)
		questions, _ := h.questionSvc.ListByGameId(ctx, lobby.GameId)
		err := h.lobbySvc.Update(context.Background(), lobby.UUID)
		if err != nil {
			log.Println("OOPS UPDATE FAIL")
		}
		h.sessions.mu.Lock()
		for _, l := range h.sessions.activeConnections[lobbyUUID] {
			l.Connection.WriteJSON(gin.H{
				"type": "questions",
				"data": questions,
			})
		}
		temp := h.sessions.activeConnections[lobbyUUID][lobbyUUID]
		temp.GameId = lobby.GameId
		temp.QuestionCount = len(questions)
		h.sessions.activeConnections[lobbyUUID][lobbyUUID] = temp
		h.sessions.mu.Unlock()
		return
	}

	if string(msg) == "end_lobby" {
		h.sessions.mu.Lock()
		for _, l := range h.sessions.activeConnections[lobbyUUID] {
			l.Connection.WriteJSON(gin.H{
				"type": "end_lobby",
			})
		}
		h.sessions.mu.Unlock()
		h.gameSvc.CalcResultNum(ctx, lobbyUUID)
		answers := h.gameSvc.GetTextAnswers(ctx, lobbyUUID)
		for _, l := range h.sessions.activeConnections[lobbyUUID] {
			l.Connection.WriteJSON(gin.H{
				"type": "text_questions",
				"data": answers,
			})
		}
		return
	}

	if string(msg) == "calculate_result" {
		data := h.gameSvc.CalculateQuizResult(ctx, lobbyUUID)
		h.sessions.mu.Lock()
		h.sessions.activeConnections[lobbyUUID][lobbyUUID].Connection.WriteJSON(gin.H{
			"type": "quiz_result",
			"data": data,
		})
		h.sessions.mu.Unlock()
		return
	}

	if strings.Contains(string(msg), "next_question:") {
		id := 0
		fmt.Sscanf(string(msg), "next_question:%d", &id)
		h.sessions.mu.Lock()
		count := h.sessions.activeConnections[lobbyUUID][lobbyUUID].QuestionCount
		for _, l := range h.sessions.activeConnections[lobbyUUID] {
			l.Connection.WriteJSON(gin.H{
				"type":           "next_question",
				"data":           id,
				"question_count": count,
			})
		}
		h.sessions.mu.Unlock()
		return
	}

	if strings.Contains(string(msg), "get_question:") {
		questionNum := 0
		isText := false
		fmt.Sscanf(string(msg), "get_question:%d", &questionNum)
		lobby, _ := h.lobbySvc.LoadByUUID(ctx, lobbyUUID)
		question, _ := h.questionSvc.LoadByNumber(ctx, lobby.GameId, questionNum)
		if question.AnswerText != "" {
			isText = true
		}
		h.sessions.mu.Lock()
		for _, l := range h.sessions.activeConnections[lobbyUUID] {
			l.Connection.WriteJSON(gin.H{
				"type":   "question",
				"data":   question,
				"isText": isText,
			})
		}
		h.sessions.mu.Unlock()
		return
	}

	if strings.Contains(string(msg), "answer_num:") {
		questionId := 0
		questionNum := 0
		answer := 0
		fmt.Sscanf(string(msg), "answer_num:%d:%d:%d", &questionId, &questionNum, &answer)

		data := model.Answer{
			LobbyUUID:      lobbyUUID,
			PlayerUUID:     playerUUID,
			AnswerNum:      answer,
			AnswerText:     "",
			QuestionNumber: questionNum,
			QuestionId:     questionId,
		}
		h.gameSvc.SaveAnswer(ctx, data)
		h.sessions.mu.Lock()
		playerName := h.sessions.activeConnections[lobbyUUID][playerUUID].UserName
		h.sessions.activeConnections[lobbyUUID][lobbyUUID].Connection.WriteJSON(gin.H{
			"type": "answer",
			"data": playerName,
		})
		h.sessions.mu.Unlock()
		return
	}

	if strings.Contains(string(msg), "answer_text:") {
		questionId, questionNum, answerText := parseTextAnswer(string(msg))
		data := model.Answer{
			LobbyUUID:      lobbyUUID,
			PlayerUUID:     playerUUID,
			AnswerNum:      0,
			AnswerText:     answerText,
			QuestionNumber: questionNum,
			QuestionId:     questionId,
		}
		h.gameSvc.SaveAnswer(ctx, data)
		h.sessions.mu.Lock()
		playerName := h.sessions.activeConnections[lobbyUUID][playerUUID].UserName
		h.sessions.activeConnections[lobbyUUID][lobbyUUID].Connection.WriteJSON(gin.H{
			"type": "answer",
			"data": playerName,
		})
		h.sessions.mu.Unlock()
		return
	}

	if strings.Contains(string(msg), "result_text:") {
		res := strings.Split(string(msg), ":")
		pUUID, _ := uuid.Parse(res[1])
		qNum, _ := strconv.Atoi(res[2])
		isCorrect := false
		if res[3] == "1" {
			isCorrect = true
		}

		h.sessions.mu.Lock()
		gameId := h.sessions.activeConnections[lobbyUUID][lobbyUUID].GameId
		h.sessions.mu.Unlock()
		result := model.SaveTextResult{
			LobbyUUID:      lobbyUUID,
			PlayerUUID:     pUUID,
			QuestionNumber: qNum,
			IsCorrect:      isCorrect,
			GameId:         gameId,
		}

		log.Println("result:", result)
		h.gameSvc.SaveTextResult(ctx, result)
	}
}

func parseTextAnswer(input string) (int, int, string) {
	parts := strings.SplitN(input, "/", 2)

	firstPart := parts[0]

	text := parts[1]

	questionId := 0
	questionNum := 0

	fmt.Sscanf(firstPart, "answer_text:%d:%d", &questionId, &questionNum)

	return questionId, questionNum, text
}
