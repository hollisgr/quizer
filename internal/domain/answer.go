package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrAnswerNotFound = errors.New("answer not found")
)

type Answer struct {
	Id             int       `json:"answer_id" db:"id"`
	LobbyUUID      uuid.UUID `json:"lobby_uuid" db:"lobby_uuid"`
	PlayerUUID     uuid.UUID `json:"player_uuid" db:"player_uuid"`
	AnswerNum      int       `json:"answer_num" db:"answer_num"`
	AnswerText     string    `json:"answer_text" db:"answer_text"`
	QuestionNumber int       `json:"question_num" db:"question_num"`
	QuestionId     int       `json:"question_id" db:"question_id"`
}

type AnswerWithData struct {
	Answer
	Description   string `json:"description" db:"description"`
	UserName      string `json:"user_name" db:"user_name"`
	CorrectAnswer string `json:"correct_answer" db:"correct_answer"`
}

type AnswerRepository interface {
	SaveAnswer(ctx context.Context, data Answer) error
	LoadAnswersByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) ([]Answer, error)
	LoadTextAnswersByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) ([]AnswerWithData, error)
	LoadTextAnswer(ctx context.Context, lobbyUUID uuid.UUID, playerUUID uuid.UUID, questionNum int) (AnswerWithData, error)
}
