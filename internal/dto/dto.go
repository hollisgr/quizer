package dto

type CreateNewGame struct {
	OwnerId     int
	Description string
	Link        string
}

type CreateNewGameRequest struct {
	Description string `json:"description"`
	Link        string `json:"link"`
}

type CreateNewQuestionRequest struct {
	GameId      int    `json:"game_id" db:"game_id"`
	Number      int    `json:"number" db:"number"`
	Cost        int    `json:"cost" db:"cost"`
	AnswerNum   int    `json:"answer" db:"answer"`
	AnswerText  string `json:"answer_text" db:"answer_text"`
	Description string `json:"description" db:"description"`
}
