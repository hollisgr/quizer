package db

import (
	"context"
	"fmt"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *storage) CreateQuestion(ctx context.Context, data dto.CreateNewQuestionRequest) (int, error) {
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

func (s *storage) QuestionsByGameId(ctx context.Context, gameId int) ([]model.Question, error) {
	var res []model.Question
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

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.Question])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) QuestionLoad(ctx context.Context, id int) (model.Question, error) {
	var res model.Question
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

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Question])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) QuestionLoadByNumber(ctx context.Context, gameId int, number int) (model.Question, error) {
	var res model.Question
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

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Question])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) UpdateQuestion(ctx context.Context, updated model.Question) (int, error) {
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

func (s *storage) DeleteQuestion(ctx context.Context, id int) (int, error) {
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
