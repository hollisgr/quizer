package domain

import (
	"context"

	"github.com/google/uuid"
)

type Result struct {
	Id             int       `json:"result_id" db:"id"`
	LobbyUUID      uuid.UUID `json:"lobby_uuid" db:"lobby_uuid"`
	PlayerUUID     uuid.UUID `json:"player_uuid" db:"player_uuid"`
	QuestionNumber int       `json:"question_num" db:"question_num"`
	QuestionId     int       `json:"question_id" db:"question_id"`
	AnswerNumber   int       `json:"answer_num" db:"answer_num"`
	AnswerText     string    `json:"answer_text" db:"answer_text"`
	Score          int       `json:"score" db:"score"`
}

type SaveTextResult struct {
	LobbyUUID      uuid.UUID `json:"lobby_uuid" db:"lobby_uuid"`
	PlayerUUID     uuid.UUID `json:"player_uuid" db:"player_uuid"`
	QuestionId     int       `json:"question_id" db:"question_id"`
	QuestionNumber int       `json:"question_num" db:"question_num"`
	IsCorrect      bool
	GameId         int
}

type CalcResult struct {
	TotalScore int    `json:"total_score" db:"total_score"`
	UserName   string `json:"user_name" db:"user_name"`
}

type ResultRepository interface {
	SaveResult(ctx context.Context, data Result) error
	LoadResultByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) []Result
	LoadPlayerResult(ctx context.Context, lobbyUUID uuid.UUID, playerUUID uuid.UUID) []Result
	CalculateResults(ctx context.Context, lobbyUUID uuid.UUID) []CalcResult
}
