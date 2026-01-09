package postgres

import (
	"context"
	"fmt"
	"quizer_server/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type resultStorage struct {
	db *pgxpool.Pool
}

func NewResultStorage(pool *pgxpool.Pool) domain.ResultRepository {
	return &resultStorage{
		db: pool,
	}
}

func (s *resultStorage) SaveResult(ctx context.Context, data domain.Result) error {
	id := 0
	query := `
		INSERT INTO
			player_results (
				lobby_uuid,
				player_uuid,
				question_num,
				question_id,
				answer_num,
				answer_text,
				score
			)
		VALUES
			(
			@lobby_uuid,
			@player_uuid,
			@question_num,
			@question_id,
			@answer_num,
			@answer_text,
			@score
		)
		RETURNING
			id
	`
	args := pgx.NamedArgs{
		"lobby_uuid":   data.LobbyUUID,
		"player_uuid":  data.PlayerUUID,
		"question_num": data.QuestionNumber,
		"question_id":  data.QuestionId,
		"answer_num":   data.AnswerNumber,
		"answer_text":  data.AnswerText,
		"score":        data.Score,
	}
	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("db save result error: %v", err)
	}
	return nil
}

func (s *resultStorage) LoadResultByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) []domain.Result {
	res := []domain.Result{}
	query := `
		SELECT
			lobby_uuid,
			player_uuid,
			question_num,
			question_id,
			answer_num,
			answer_text,
			score
		FROM player_results
		WHERE lobby_uuid = @lobby_uuid
		ORDER BY id desc
	`
	args := pgx.NamedArgs{
		"lobby_uuid": lobbyUUID,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res
	}

	res, _ = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Result])

	return res
}

func (s *resultStorage) LoadPlayerResult(ctx context.Context, lobbyUUID uuid.UUID, playerUUID uuid.UUID) []domain.Result {
	res := []domain.Result{}
	query := `
		SELECT
			lobby_uuid,
			player_uuid,
			question_num,
			question_id,
			answer_num,
			answer_text,
			score
		FROM player_results
		WHERE 
			lobby_uuid = @lobby_uuid
			AND
			player_uuid = @player_uuid
		ORDER BY id desc
	`
	args := pgx.NamedArgs{
		"lobby_uuid":  lobbyUUID,
		"player_uuid": playerUUID,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res
	}

	res, _ = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Result])

	return res
}

func (s *resultStorage) CalculateResults(ctx context.Context, lobbyUUID uuid.UUID) []domain.CalcResult {
	res := []domain.CalcResult{}
	query := `
		SELECT 
			p.user_name, 
			SUM(score) as total_score 
		FROM player_results pa 
		JOIN players p ON p.uuid = pa.player_uuid
		WHERE lobby_uuid = @lobby_uuid
		GROUP BY p.user_name 
		ORDER BY total_score DESC
	`
	args := pgx.NamedArgs{
		"lobby_uuid": lobbyUUID,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res
	}

	res, _ = pgx.CollectRows(rows, pgx.RowToStructByName[domain.CalcResult])

	return res
}
