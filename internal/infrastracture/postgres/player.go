package postgres

import (
	"context"
	"fmt"
	"quizer_server/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type playerStorage struct {
	db *pgxpool.Pool
}

func NewPlayerStorage(pool *pgxpool.Pool) domain.PlayerRepository {
	return &playerStorage{
		db: pool,
	}
}

func (s *playerStorage) SavePlayer(ctx context.Context, newPlayer domain.Player) error {
	var uuid uuid.UUID
	query := `
		INSERT INTO
			players (
				uuid,
				lobby_id,
				user_name,
				is_admin
			)
		VALUES
			(
			@uuid,
			@lobby_id,
			@user_name,
			@is_admin
		)
		RETURNING
			uuid
	`
	args := pgx.NamedArgs{
		"uuid":      newPlayer.UUID,
		"lobby_id":  newPlayer.LobbyUUID,
		"user_name": newPlayer.UserName,
		"is_admin":  newPlayer.IsAdmin,
	}
	row := s.db.QueryRow(ctx, query, args)
	err := row.Scan(&uuid)
	if err != nil {
		return fmt.Errorf("db create new game error: %v", err)
	}
	return nil
}

func (s *playerStorage) PlayersByGameUUID(ctx context.Context, gameUUID uuid.UUID) ([]domain.Player, error) {
	var res []domain.Player
	query := `
		SELECT
			uuid,
			lobby_id,
			user_name,
			is_admin
		FROM players
		WHERE uuid = @gameUUID
		ORDER BY id desc
	`
	args := pgx.NamedArgs{
		"gameUUID": gameUUID,
	}
	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Player])

	if err != nil {
		return res, err
	}

	return res, nil
}
