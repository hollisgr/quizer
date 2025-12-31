package db

import (
	"context"
	"fmt"
	"log"
	"quizer_server/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *storage) CreateLobby(ctx context.Context, data model.Lobby) error {
	var id uuid.UUID
	query := `
		INSERT INTO
			lobbies (
				uuid,
				game_id,
				is_started
			)
		VALUES
			(
			@uuid,
			@game_id,
			@is_started
		)
		RETURNING
			uuid
	`
	args := pgx.NamedArgs{
		"uuid":       data.UUID,
		"game_id":    data.GameId,
		"is_started": data.IsStarted,
	}
	err := s.db.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return fmt.Errorf("db create new lobby error: %v", err)
	}
	return nil
}

func (s *storage) LobbyLoadByUUID(ctx context.Context, uuid uuid.UUID) (model.Lobby, error) {
	var res model.Lobby
	query := `
		SELECT
			uuid,
			game_id,
			is_started
		FROM lobbies 
		WHERE uuid = @uuid
	`

	args := pgx.NamedArgs{
		"uuid": uuid,
	}

	rows, err := s.db.Query(ctx, query, args)
	defer rows.Close()

	if err != nil {
		return res, err
	}

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Lobby])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) LobbyList(ctx context.Context) ([]model.Lobby, error) {
	res := []model.Lobby{}
	query := `
		SELECT
			uuid,
			game_id,
			is_started
		FROM lobbies
		WHERE is_started = false
	`
	rows, err := s.db.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		return res, err
	}
	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.Lobby])
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *storage) UpdateLobby(ctx context.Context, lobbyUUID uuid.UUID) error {
	log.Println("db update, uuid:", lobbyUUID)
	query := `
		UPDATE
			lobbies
		SET
			is_started = true
		WHERE uuid = @lobbyUUID
	`
	args := pgx.NamedArgs{
		"lobbyUUID": lobbyUUID,
	}
	_, err := s.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db update lobby error: %v", err)
	}
	return nil
}
