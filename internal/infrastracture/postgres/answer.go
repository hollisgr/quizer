package postgres

import (
	"context"
	"errors"
	"fmt"
	"quizer_server/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type answerStorage struct {
	db *pgxpool.Pool
}

func NewAnswerStorage(pool *pgxpool.Pool) domain.AnswerRepository {
	return &answerStorage{
		db: pool,
	}
}

func (s *answerStorage) SaveAnswer(ctx context.Context, data domain.Answer) error {
	id := 0
	query := `
		INSERT INTO
			player_answers (
				lobby_uuid,
				player_uuid,
				answer_num,
				answer_text,
				question_num,
				question_id
			)
		VALUES
			(
			@lobby_uuid,
			@player_uuid,
			@answer_num,
			@answer_text,
			@question_num,
			@question_id
		)
		RETURNING
			id
	`
	args := pgx.NamedArgs{
		"lobby_uuid":   data.LobbyUUID,
		"player_uuid":  data.PlayerUUID,
		"answer_num":   data.AnswerNum,
		"answer_text":  data.AnswerText,
		"question_num": data.QuestionNumber,
		"question_id":  data.QuestionId,
	}
	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("db create new lobby error: %v", err)
	}
	return nil
}

func (s *answerStorage) LoadAnswersByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) ([]domain.Answer, error) {
	res := []domain.Answer{}
	query := `
		SELECT
			id,
			lobby_uuid,
			player_uuid,
			answer_num,
			answer_text,
			question_num,
			question_id
		FROM player_answers
		WHERE lobby_uuid = @lobby_uuid
		ORDER BY id desc
	`
	args := pgx.NamedArgs{
		"lobby_uuid": lobbyUUID,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Answer])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res, domain.ErrAnswerNotFound
		}
		return res, err
	}

	return res, nil
}

func (s *answerStorage) LoadTextAnswersByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) ([]domain.AnswerWithData, error) {
	res := []domain.AnswerWithData{}
	query := `
		SELECT
			pa.lobby_uuid,
			pa.player_uuid,
			p.user_name,
			description,
			pa.answer_text,
			q.answer_text AS correct_answer,
			pa.question_num,
			pa.question_id
		FROM player_answers pa
		JOIN questions q ON pa.question_id = q.id
		JOIN players p ON pa.player_uuid = p.uuid 
		WHERE pa.lobby_uuid = @lobby_uuid 
		AND pa.answer_text != ''
		ORDER BY pa.id ASC;
	`
	args := pgx.NamedArgs{
		"lobby_uuid": lobbyUUID,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.AnswerWithData])

	if err != nil {
		return res, err
	}

	if len(res) == 0 {
		return res, domain.ErrAnswerNotFound
	}

	return res, nil
}

func (s *answerStorage) LoadTextAnswer(ctx context.Context, lobbyUUID uuid.UUID, playerUUID uuid.UUID, questionNum int) (domain.AnswerWithData, error) {
	res := domain.AnswerWithData{}
	query := `
		SELECT
			pa.lobby_uuid,
			pa.player_uuid,
			p.user_name,
			description,
			pa.answer_text,
			q.answer_text AS correct_answer,
			pa.question_num,
			pa.question_id
		FROM player_answers pa
		JOIN questions q ON pa.question_id = q.id
		JOIN players p ON pa.player_uuid = p.uuid 
		WHERE 
			pa.lobby_uuid = @lobby_uuid
			AND pa.player_uuid = @player_uuid
			AND pa.question_num = @question_num
			AND pa.answer_text != ''
	`
	args := pgx.NamedArgs{
		"lobby_uuid":   lobbyUUID,
		"player_uuid":  playerUUID,
		"question_num": questionNum,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.AnswerWithData])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res, domain.ErrAnswerNotFound
		}
		return res, err
	}

	return res, nil
}
