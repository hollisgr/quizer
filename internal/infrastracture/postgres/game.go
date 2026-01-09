package postgres

import (
	"context"
	"fmt"
	"quizer_server/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type gameStorage struct {
	db *pgxpool.Pool
}

func NewGameStorage(pool *pgxpool.Pool) domain.GameRepository {
	return &gameStorage{
		db: pool,
	}
}

func (s *gameStorage) CreateGame(ctx context.Context, data domain.Game) (int, error) {
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
		"owner_id":    data.Owner,
		"link":        data.Link,
	}
	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("db create new game error: %v", err)
	}
	return id, nil
}

func (s *gameStorage) GameList(ctx context.Context) ([]domain.GameWithLogin, error) {
	var res []domain.GameWithLogin
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

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.GameWithLogin])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *gameStorage) GameLoad(ctx context.Context, id int) (domain.GameWithLogin, error) {
	var res domain.GameWithLogin
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

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.GameWithLogin])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *gameStorage) UpdateGame(ctx context.Context, updated domain.Game) (int, error) {
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

func (s *gameStorage) UpdateFilePath(ctx context.Context, gameId int, path string) (int, error) {
	res := 0
	query := `
		UPDATE
			games
		SET
			link = @link
		WHERE
			id = @id
		RETURNING id
	`
	args := pgx.NamedArgs{
		"id":   gameId,
		"link": path,
	}
	row := s.db.QueryRow(ctx, query, args)

	err := row.Scan(&res)

	if err != nil || res == 0 {
		return res, err
	}

	return res, nil
}

func (s *gameStorage) DeleteGame(ctx context.Context, id int) (int, error) {
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
