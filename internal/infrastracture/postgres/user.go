package postgres

import (
	"context"
	"errors"
	"fmt"
	"quizer_server/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userStorage struct {
	db *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) domain.UserRepository {
	return &userStorage{
		db: pool,
	}
}

func (s *userStorage) ByLogin(ctx context.Context, login string) (domain.User, error) {
	var res domain.User
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
		return res, fmt.Errorf("db user by login query exec err: %w", err)
	}

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[domain.User])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return res, domain.ErrUserNotFound
		}
		return res, fmt.Errorf("db user by login collecting row err: %w", err)
	}

	return res, nil
}
