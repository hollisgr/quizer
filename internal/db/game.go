package db

import (
	"context"
	"fmt"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"

	"github.com/google/uuid"
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

func (s *storage) SavePlayer(ctx context.Context, newPlayer model.Player) error {
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

func (s *storage) PlayersByGameUUID(ctx context.Context, gameUUID uuid.UUID) ([]model.Player, error) {
	var res []model.Player
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

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.Player])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) SaveAnswer(ctx context.Context, data model.Answer) error {
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

func (s *storage) LoadAnswersByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) ([]model.Answer, error) {
	res := []model.Answer{}
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

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.Answer])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) LoadTextAnswersByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) ([]model.PlayerTextAnswer, error) {
	res := []model.PlayerTextAnswer{}
	query := `
		SELECT
			pa.lobby_uuid,
			pa.player_uuid,
			p.user_name,
			description,
			pa.answer_text AS player_answer,
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

	res, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.PlayerTextAnswer])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) LoadTextAnswer(ctx context.Context, lobbyUUID uuid.UUID, playerUUID uuid.UUID, questionNum int) (model.PlayerTextAnswer, error) {
	res := model.PlayerTextAnswer{}
	query := `
		SELECT
			pa.lobby_uuid,
			pa.player_uuid,
			p.user_name,
			description,
			pa.answer_text AS player_answer,
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

	res, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.PlayerTextAnswer])

	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *storage) SaveResult(ctx context.Context, data model.Result) error {
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

func (s *storage) LoadResultByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) []model.Result {
	res := []model.Result{}
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

	res, _ = pgx.CollectRows(rows, pgx.RowToStructByName[model.Result])

	return res
}

func (s *storage) LoadPlayerResult(ctx context.Context, lobbyUUID uuid.UUID, playerUUID uuid.UUID) []model.Result {
	res := []model.Result{}
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

	res, _ = pgx.CollectRows(rows, pgx.RowToStructByName[model.Result])

	return res
}

func (s *storage) CalculateResults(ctx context.Context, lobbyUUID uuid.UUID) []model.CalcResult {
	res := []model.CalcResult{}
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

	res, _ = pgx.CollectRows(rows, pgx.RowToStructByName[model.CalcResult])

	return res
}
