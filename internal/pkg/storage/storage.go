package storage

import (
	"monotonic/internal/pkg/models"
)

var (
	userWords     = make(map[int64][]int)
	userQuestions = make(map[int64]models.PracticeQuestion)
)

func AddWord(userID int64, wordID int) {
	userWords[userID] = append(userWords[userID], wordID)
}

func ClearList(userID int64) {
	if _, ok := userWords[userID]; ok {
		delete(userWords, userID)
	}
}

func GetUserWords(userID int64) ([]int, bool) {
	words, ok := userWords[userID]
	return words, ok
}

func SaveQuestion(userID int64, q models.PracticeQuestion) {
	userQuestions[userID] = q
}

func GetQuestion(userID int64) (models.PracticeQuestion, bool) {
	q, ok := userQuestions[userID]
	return q, ok
}

func ClearQuestion(userID int64) {
	delete(userQuestions, userID)
}
