package db

import (
	"context"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	UserByLogin(ctx context.Context, login string) (model.User, error)

	CreateGame(ctx context.Context, data dto.CreateNewGame) (int, error)
	GameList(ctx context.Context) ([]model.Game, error)
	GameLoad(ctx context.Context, id int) (model.Game, error)
	UpdateGame(ctx context.Context, updated model.Game) (int, error)
	DeleteGame(ctx context.Context, id int) (int, error)

	CreateLobby(ctx context.Context, data model.Lobby) error
	LobbyLoadByUUID(ctx context.Context, uuid uuid.UUID) (model.Lobby, error)

	PlayersByGameUUID(ctx context.Context, gameUUID uuid.UUID) ([]model.Player, error)
	SavePlayer(ctx context.Context, newPlayer model.Player) error
	SaveAnswer(ctx context.Context, data model.Answer) error

	LoadAnswersByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) ([]model.Answer, error)
	SaveResult(ctx context.Context, data model.Result) error
	LoadResultByLobbyUUID(ctx context.Context, lobbyUUID uuid.UUID) []model.Result
	LoadPlayerResult(ctx context.Context, lobbyUUID uuid.UUID, playerUUID uuid.UUID) []model.Result

	CreateQuestion(ctx context.Context, data dto.CreateNewQuestionRequest) (int, error)
	QuestionLoad(ctx context.Context, id int) (model.Question, error)
	QuestionLoadByNumber(ctx context.Context, gameId int, number int) (model.Question, error)

	QuestionsByGameId(ctx context.Context, gameId int) ([]model.Question, error)
	UpdateQuestion(ctx context.Context, updated model.Question) (int, error)
	DeleteQuestion(ctx context.Context, id int) (int, error)
}

type storage struct {
	db *pgxpool.Pool
}

func New(p *pgxpool.Pool) Storage {
	return &storage{
		db: p,
	}
}
