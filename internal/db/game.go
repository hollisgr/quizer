package db

import (
	"context"
	"fmt"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *storage) CreateGame(ctx context.Context, data dto.CreateNewGame) (int, error) {
	var id int
	query := `
		INSERT INTO
			games (
				description,
				owner_id,
				link
			)
		VALUES
			(
			@description,
			@owner_id,
			@link
		)
		RETURNING
			id
	`
	args := pgx.NamedArgs{
		"description": data.Description,
		"owner_id":    data.OwnerId,
		"link":        data.Link,
	}
	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("db create new game error: %v", err)
	}
	return id, nil
}

func (s *storage) GameList(ctx context.Context) ([]model.Game, error) {
	var res []model.Game
	query := `
		SELECT
			g.id, 
			description, 
			login, 
			created_at, 
			link
		FROM games g 
		JOIN users u on u.id = g.owner_id
		ORDER BY id desc
	`
	rows, err := s.db.Query(ctx, query)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.Game])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) GameLoad(ctx context.Context, id int) (model.Game, error) {
	var res model.Game
	query := `
		SELECT
			g.id, 
			description, 
			login, 
			created_at, 
			link
		FROM games g 
		JOIN users u on u.id = g.owner_id
		WHERE g.id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Game])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) UpdateGame(ctx context.Context, updated model.Game) (int, error) {
	res := 0
	query := `
		UPDATE
			games
		SET
			description = @description,
			link = @link
		WHERE
			id = @id
		RETURNING id
	`
	args := pgx.NamedArgs{
		"id":          updated.Id,
		"description": updated.Description,
		"link":        updated.Link,
	}
	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&res)

	if err != nil || res == 0 {
		return res, err
	}

	return res, nil
}

func (s *storage) DeleteGame(ctx context.Context, id int) (int, error) {
	res := 0
	query := `
		DELETE FROM
			games
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