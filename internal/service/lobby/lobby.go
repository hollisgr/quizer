package lobby

import (
	"context"
	"log"
	"quizer_server/internal/db"
	"quizer_server/internal/model"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, lobby model.Lobby) (int, error)
	LoadByUUID(ctx context.Context, uuid uuid.UUID) (model.Lobby, error)
	List(ctx context.Context) ([]model.Lobby, error)
	Update(ctx context.Context, lobbyUUID uuid.UUID) error
}

type lobbyService struct {
	storage db.Storage
}

func New(s db.Storage) Service {
	return &lobbyService{
		storage: s,
	}
}

func (ls *lobbyService) Create(ctx context.Context, lobby model.Lobby) (int, error) {
	count := 0
	lobby.IsStarted = false
	err := ls.storage.CreateLobby(ctx, lobby)

	if err != nil {
		log.Println("lobby svc create err:", err)
	}

	questions, _ := ls.storage.QuestionsByGameId(ctx, lobby.GameId)

	count = len(questions)

	return count, nil
}

func (ls *lobbyService) LoadByUUID(ctx context.Context, uuid uuid.UUID) (model.Lobby, error) {
	res, err := ls.storage.LobbyLoadByUUID(ctx, uuid)
	if err != nil {
		log.Println("lobby svc load err", err)
		return res, err
	}
	return res, nil
}

func (ls *lobbyService) List(ctx context.Context) ([]model.Lobby, error) {
	res, err := ls.storage.LobbyList(ctx)
	if err != nil {
		log.Println("lobby svc list err:", err)
		return res, err
	}
	return res, nil
}

func (ls *lobbyService) Update(ctx context.Context, lobbyUUID uuid.UUID) error {
	log.Println("svc update, uuid:", lobbyUUID)
	err := ls.storage.UpdateLobby(ctx, lobbyUUID)
	if err != nil {
		log.Println("lobby svc update err:", err)
		return err
	}
	return nil
}
