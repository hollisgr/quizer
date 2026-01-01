package game

import (
	"context"
	"log"
	"quizer_server/internal/db"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"

	"github.com/google/uuid"
)

type Service interface {
	CreateNewGame(ctx context.Context, data dto.CreateNewGame) (int, error)
	GameList(ctx context.Context) ([]model.Game, error)
	GameLoad(ctx context.Context, id int) (model.Game, error)
	DeleteGame(ctx context.Context, id int) (int, error)
	UpdateGame(ctx context.Context, updated model.Game) (int, error)
	UpdateFilePath(ctx context.Context, gameId int, path string) (int, error)

	GetPlayersByGameUUID(ctx context.Context, gameUUID uuid.UUID) []model.Player
	SavePlayer(ctx context.Context, newPlayer model.Player) error

	SaveAnswer(ctx context.Context, data model.Answer)
	GetTextAnswers(ctx context.Context, lobbyUUID uuid.UUID) []model.PlayerTextAnswer

	CalcResultNum(ctx context.Context, lobbyUUID uuid.UUID)
	CalculateQuizResult(ctx context.Context, lobbyUUID uuid.UUID) []model.CalcResult
	SaveTextResult(ctx context.Context, result model.SaveTextResult)
}

type gameService struct {
	storage db.Storage
}

func New(s db.Storage) Service {
	return &gameService{
		storage: s,
	}
}

func (gs *gameService) CreateNewGame(ctx context.Context, data dto.CreateNewGame) (int, error) {
	id, err := gs.storage.CreateGame(ctx, data)
	if err != nil {
		log.Println(err)
		return id, err
	}
	return id, err
}

func (gs *gameService) GameList(ctx context.Context) ([]model.Game, error) {
	list, err := gs.storage.GameList(ctx)
	if err != nil {
		log.Println(err)
		return list, err
	}
	return list, nil
}

func (gs *gameService) GameLoad(ctx context.Context, id int) (model.Game, error) {
	res, err := gs.storage.GameLoad(ctx, id)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, nil
}

func (gs *gameService) UpdateGame(ctx context.Context, updated model.Game) (int, error) {
	res, err := gs.storage.UpdateGame(ctx, updated)
	if err != nil || res == 0 {
		log.Println(err)
		return res, err
	}
	return res, nil
}

func (gs *gameService) DeleteGame(ctx context.Context, id int) (int, error) {
	res, err := gs.storage.DeleteGame(ctx, id)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, nil
}

func (gs *gameService) SavePlayer(ctx context.Context, newPlayer model.Player) error {
	gs.storage.SavePlayer(ctx, newPlayer)
	return nil
}

func (gs *gameService) GetPlayersByGameUUID(ctx context.Context, gameUUID uuid.UUID) []model.Player {
	res, _ := gs.storage.PlayersByGameUUID(ctx, gameUUID)
	return res
}

func (gs *gameService) SaveAnswer(ctx context.Context, data model.Answer) {
	err := gs.storage.SaveAnswer(ctx, data)
	if err != nil {
		log.Println("service save answer err: ", err)
	}
}

func (gs *gameService) GetTextAnswers(ctx context.Context, lobbyUUID uuid.UUID) []model.PlayerTextAnswer {
	answers, err := gs.storage.LoadTextAnswersByLobbyUUID(ctx, lobbyUUID)
	if err != nil {
		log.Println("game service get text answers err:", err)
		return answers
	}
	return answers
}

func (gs *gameService) CalcResultNum(ctx context.Context, lobbyUUID uuid.UUID) {
	log.Println("calcucating result")
	lobby, err := gs.storage.LobbyLoadByUUID(ctx, lobbyUUID)
	if err != nil {
		log.Println("calc result num load lobby err: ", err)
	}

	answers, err := gs.storage.LoadAnswersByLobbyUUID(ctx, lobbyUUID)
	if err != nil {
		log.Println("calc result num load answers err: ", err)
	}

	if len(answers) == 0 {
		log.Println("calc result num answer list is empty")
	}

	qArr, err := gs.storage.QuestionsByGameId(ctx, lobby.GameId)
	if err != nil {
		log.Println("calc result num load questions err: ", err)

	}

	for _, a := range answers {
		for _, q := range qArr {
			if a.AnswerNum != 0 && a.QuestionNumber == q.Number {
				score := 0
				if a.AnswerNum == q.AnswerNum {
					score = q.Cost
				}
				res := model.Result{
					LobbyUUID:      lobbyUUID,
					PlayerUUID:     a.PlayerUUID,
					QuestionNumber: a.QuestionNumber,
					QuestionId:     a.QuestionId,
					AnswerNumber:   a.AnswerNum,
					Score:          score,
				}
				err = gs.storage.SaveResult(ctx, res)
				if err != nil {
					log.Println("calc result num save result err: ", err)

				}
			}
		}
	}
}

func (gs *gameService) SaveTextResult(ctx context.Context, data model.SaveTextResult) {
	question, err := gs.storage.QuestionLoadByNumber(ctx, data.GameId, data.QuestionNumber)
	if err != nil {
		log.Println("game service save text result question load err:", err)
		return
	}
	answer, err := gs.storage.LoadTextAnswer(ctx, data.LobbyUUID, data.PlayerUUID, data.QuestionNumber)
	if err != nil {
		log.Println("game service save text result answer load err:", err)
		return
	}
	score := 0
	if data.IsCorrect {
		score = 1
	}
	result := model.Result{
		LobbyUUID:      data.LobbyUUID,
		PlayerUUID:     data.PlayerUUID,
		QuestionNumber: data.QuestionNumber,
		QuestionId:     answer.QuestionId,
		AnswerText:     answer.AnswerText,
		Score:          question.Cost * score,
	}
	err = gs.storage.SaveResult(ctx, result)
	if err != nil {
		log.Println("game service save text result save result err:", err)
		return
	}
}

func (gs *gameService) CalculateQuizResult(ctx context.Context, lobbyUUID uuid.UUID) []model.CalcResult {
	res := gs.storage.CalculateResults(ctx, lobbyUUID)
	return res
}

func (gs *gameService) UpdateFilePath(ctx context.Context, gameId int, path string) (int, error) {
	id, err := gs.storage.UpdateFilePath(ctx, gameId, path)
	if err != nil {
		log.Println("update file path err:", err)
		return id, err
	}
	return id, nil
}
