package postgres

import (
	"context"
	"fmt"
	"quizer_server/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type questionStorage struct {
	db *pgxpool.Pool
}

func NewQuestionStorage(pool *pgxpool.Pool) domain.QuestionRepository {
	return &questionStorage{
		db: pool,
	}
}

func (s *questionStorage) CreateQuestion(ctx context.Context, data domain.Question) (int, error) {
	var id int
	query := `
		INSERT INTO
			questions (
 				number,
 				description,
 				game_id,
 				answer,
 				answer_text,
 				cost
			)
		VALUES
			(
 			@number,
 			@description,
 			@game_id,
 			@answer,
 			@answer_text,
 			@cost
		)
		RETURNING
			id
	`
	args := pgx.NamedArgs{
		"number":      data.Number,
		"description": data.Description,
		"game_id":     data.GameId,
		"answer":      data.AnswerNum,
		"answer_text": data.AnswerText,
		"cost":        data.Cost,
	}
	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("db create new game error: %v", err)
	}
	return id, nil
}

func (s *questionStorage) QuestionsByGameId(ctx context.Context, gameId int) ([]domain.Question, error) {
	var res []domain.Question
	query := `
		SELECT
			id,
			number,
			description,
			game_id,
			answer,
			answer_text,
			cost
		FROM questions
		WHERE
			game_id = @game_id
		ORDER BY number
	`

	args := pgx.NamedArgs{
		"game_id": gameId,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Question])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *questionStorage) QuestionLoad(ctx context.Context, id int) (domain.Question, error) {
	var res domain.Question
	query := `
		SELECT
			id,
			number,
			description,
			game_id,
			answer,
			answer_text,
			cost
		FROM questions
		WHERE
			id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Question])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *questionStorage) QuestionLoadByNumber(ctx context.Context, gameId int, number int) (domain.Question, error) {
	var res domain.Question
	query := `
		SELECT
			id,
			number,
			description,
			game_id,
			answer,
			answer_text,
			cost
		FROM questions
		WHERE
			game_id = @game_id
			AND
			number = @number
	`

	args := pgx.NamedArgs{
		"game_id": gameId,
		"number":  number,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.Question])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *questionStorage) UpdateQuestion(ctx context.Context, updated domain.Question) (int, error) {
	res := 0
	query := `
		UPDATE
			questions
		SET
			number = @number,
			description = @description,
			game_id = @game_id,
			answer = @answer,
			answer_text = @answer_text,
			cost = @cost
		WHERE
			id = @id
		RETURNING id
	`
	args := pgx.NamedArgs{
		"id":          updated.Id,
		"number":      updated.Number,
		"description": updated.Description,
		"game_id":     updated.GameId,
		"answer":      updated.AnswerNum,
		"answer_text": updated.AnswerText,
		"cost":        updated.Cost,
	}
	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&res)

	if err != nil || res == 0 {
		return res, err
	}

	return res, nil
}

func (s *questionStorage) DeleteQuestion(ctx context.Context, id int) (int, error) {
	res := 0
	query := `
		DELETE FROM
			questions
		WHERE
			id = @id
		RETURNING id
	`
	args := pgx.NamedArgs{
		"id": id,
	}
	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&res)

	if err != nil || res == 0 {
		return res, err
	}

	return res, nil
}
