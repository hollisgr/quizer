package domain

import "context"

type Question struct {
	Id          int    `json:"question_id" db:"id"`
	GameId      int    `json:"game_id" db:"game_id"`
	Number      int    `json:"number" db:"number"`
	Cost        int    `json:"cost" db:"cost"`
	AnswerNum   int    `json:"answer" db:"answer"`
	AnswerText  string `json:"answer_text" db:"answer_text"`
	Description string `json:"description" db:"description"`
}

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, data Question) (int, error)
	QuestionLoad(ctx context.Context, id int) (Question, error)
	QuestionLoadByNumber(ctx context.Context, gameId int, number int) (Question, error)

	QuestionsByGameId(ctx context.Context, gameId int) ([]Question, error)
	UpdateQuestion(ctx context.Context, updated Question) (int, error)
	DeleteQuestion(ctx context.Context, id int) (int, error)
}
