package usecases

import "quizer_server/internal/domain"

type QuestionUseCases struct {
	questionRepo domain.QuestionRepository
}

func NewQuestionUseCases(questionRepo domain.QuestionRepository) *QuestionUseCases {
	return &QuestionUseCases{
		questionRepo: questionRepo,
	}
}
