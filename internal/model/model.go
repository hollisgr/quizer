package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id       int    `json:"user_id" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

type Game struct {
	Id          int       `json:"game_id" db:"id"`
	Description string    `json:"description" db:"description"`
	Owner       string    `json:"owner" db:"login"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Link        string    `json:"link" db:"link"`
}

type Question struct {
	Id          int    `json:"question_id" db:"id"`
	GameId      int    `json:"game_id" db:"game_id"`
	Number      int    `json:"number" db:"number"`
	Cost        int    `json:"cost" db:"cost"`
	AnswerNum   int    `json:"answer" db:"answer"`
	AnswerText  string `json:"answer_text" db:"answer_text"`
	Description string `json:"description" db:"description"`
}

type Lobby struct {
	UUID    uuid.UUID `json:"uuid" db:"uuid"`
	Creator uuid.UUID `json:"creator_uuid" db:"creator_uuid"`
	GameId  int       `json:"game_id" db:"game_id"`
}

type Player struct {
	UUID      uuid.UUID `json:"uuid" db:"uuid"`
	UserName  string    `json:"user_name" db:"user_name"`
	LobbyUUID uuid.UUID `json:"lobby_uuid" db:"lobby_uuid"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	GameId    int       `json:"game_id" db:"game_id"`
}

type Answer struct {
	Id             int       `json:"answer_id" db:"id"`
	LobbyUUID      uuid.UUID `json:"lobby_uuid" db:"lobby_uuid"`
	PlayerUUID     uuid.UUID `json:"player_uuid" db:"player_uuid"`
	AnswerNum      int       `json:"answer_num" db:"answer_num"`
	AnswerText     string    `json:"answer_text" db:"answer_text"`
	QuestionNumber int       `json:"question_num" db:"question_num"`
}

type Result struct {
	Id             int       `json:"result_id" db:"id"`
	LobbyUUID      uuid.UUID `json:"lobby_uuid" db:"lobby_uuid"`
	PlayerUUID     uuid.UUID `json:"player_uuid" db:"player_uuid"`
	QuestionNumber int       `json:"question_num" db:"question_num"`
	AnswerNumber   int       `json:"answer_num" db:"answer_num"`
	AnswerText     string    `json:"answer_text" db:"answer_text"`
	Score          int       `json:"score" db:"score"`
}

type JwtResponce struct {
	AccessToken string `json:"access_token"`
	UserID      int    `json:"user_id"`
}

type JwtRequest struct {
	Login    string
	Password string
}
