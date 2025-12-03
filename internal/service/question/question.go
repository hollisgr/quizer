package question

import (
	"context"
	"log"
	"quizer_server/internal/db"
	"quizer_server/internal/dto"
	"quizer_server/internal/model"
)

type Service interface {
	Create(ctx context.Context, data dto.CreateNewQuestionRequest) (int, error)
	Load(ctx context.Context, id int) (model.Question, error)
	ListByGameId(ctx context.Context, gameId int) ([]model.Question, error)
	DeleteById(ctx context.Context, id int) (int, error)
	Update(ctx context.Context, data model.Question) (int, error)
}

type questionService struct {
	storage db.Storage
}

func New(s db.Storage) Service {
	return &questionService{
		storage: s,
	}
}

func (s *questionService) Create(ctx context.Context, data dto.CreateNewQuestionRequest) (int, error) {
	id, err := s.storage.CreateQuestion(ctx, data)
	if err != nil {
		log.Println(err)
		return id, err
	}
	return id, err
}

func (s *questionService) Load(ctx context.Context, id int) (model.Question, error) {
	res, err := s.storage.QuestionLoad(ctx, id)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, err
}

func (s *questionService) ListByGameId(ctx context.Context, gameId int) ([]model.Question, error) {
	res, err := s.storage.QuestionsByGameId(ctx, gameId)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, err
}

func (s *questionService) DeleteById(ctx context.Context, id int) (int, error) {
	res, err := s.storage.DeleteQuestion(ctx, id)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, err
}

func (s *questionService) Update(ctx context.Context, data model.Question) (int, error) {
	id, err := s.storage.UpdateQuestion(ctx, data)
	if err != nil {
		log.Println(err)
		return id, err
	}
	return id, err
}
