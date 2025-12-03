package model

import "time"

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

type JwtResponce struct {
	AccessToken string `json:"access_token"`
	UserID      int    `json:"user_id"`
}

type JwtRequest struct {
	Login    string
	Password string
}
