package storage

import (
	"errors"
	"log/slog"
	"monotonic/internal/pkg/models"
	"monotonic/internal/pkg/translation"
	"sync"
	"time"

	"math/rand"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type Storage struct {
	userQuestions map[int64]models.PracticeQuestion
	userWords     map[int64][]int

	qsMu    sync.Mutex
	wordsMu sync.Mutex
}

func New() *Storage {
	return &Storage{
		userQuestions: make(map[int64]models.PracticeQuestion),
		userWords:     make(map[int64][]int),
	}
}

func (s *Storage) AddWord(userID int64, wordID int) {
	s.wordsMu.Lock()
	defer s.wordsMu.Unlock()

	slog.Debug("add word triggered")
	if _, ok := s.userWords[userID]; !ok {
		s.userWords[userID] = []int{wordID}
		slog.Info("words array", slog.Any("words", s.userWords[userID]), slog.Any("user_id", userID), slog.Any("word_id", wordID))
		return
	}
	words := s.userWords[userID]
	words = append(words, wordID)
	s.userWords[userID] = words
	slog.Info("words array", slog.Any("words", s.userWords[userID]), slog.Any("user_id", userID), slog.Any("word_id", wordID))
}

func (s *Storage) ClearList(userID int64) {
	s.wordsMu.Lock()
	defer s.wordsMu.Unlock()

	if _, ok := s.userWords[userID]; ok {
		delete(s.userWords, userID)
	}
}

func (s *Storage) GetUserWords(userID int64) ([]int, bool) {
	s.wordsMu.Lock()
	defer s.wordsMu.Unlock()

	slog.Info("words array", slog.Any("words", s.userWords[userID]))
	words, ok := s.userWords[userID]
	return words, ok
}

func (s *Storage) SaveQuestion(userID int64, q models.PracticeQuestion) {
	s.qsMu.Lock()
	defer s.qsMu.Unlock()

	s.userQuestions[userID] = q
}

func (s *Storage) GetQuestion(userID int64) (models.PracticeQuestion, bool) {
	s.qsMu.Lock()
	defer s.qsMu.Unlock()

	q, ok := s.userQuestions[userID]
	return q, ok
}

func (s *Storage) ClearQuestion(userID int64) {
	s.qsMu.Lock()
	defer s.qsMu.Unlock()

	delete(s.userQuestions, userID)
}

func (s *Storage) GeneratePracticeQuestion(userID int64) (models.PracticeQuestion, error) {
	userWordIDs, ok := s.GetUserWords(userID)
	if !ok || len(userWordIDs) == 0 {
		return models.PracticeQuestion{}, errors.New("no words available for practice")
	}

	correctIndex := r.Intn(len(userWordIDs))
	correctID := userWordIDs[correctIndex]
	correctWord := translation.WordByID[correctID]

	var distractors []models.Word
	for _, w := range translation.Words {
		if w.ID != correctID {
			distractors = append(distractors, w)
		}
	}
	r.Shuffle(len(distractors), func(i, j int) {
		distractors[i], distractors[j] = distractors[j], distractors[i]
	})

	options := append([]models.Word{correctWord}, distractors[:3]...)
	r.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	q := models.PracticeQuestion{
		English:   correctWord.English,
		CorrectID: correctID,
		Options:   options,
	}

	s.SaveQuestion(userID, q)

	return q, nil
}

func (s *Storage) IsCorrectAnswer(userID int64, wordID int) bool {
	q, ok := s.GetQuestion(userID)
	if !ok {
		return false
	}

	s.ClearQuestion(userID)
	return wordID == q.CorrectID
}
