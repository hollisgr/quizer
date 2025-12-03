package db

import (
	"context"
	"quizer_server/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *storage) UserByLogin(ctx context.Context, login string) (model.User, error) {
	var res model.User
	query := `
		SELECT
			id,
			login,
			password
		FROM
			users
		WHERE
			login = @login
	`
	args := pgx.NamedArgs{
		"login": login,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])

	if err != nil {
		return res, err
	}

	return res, nil
}
