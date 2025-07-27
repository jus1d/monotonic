package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
	"monotonic/internal/pkg/translation"
	"monotonic/internal/storage/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
	r  *rand.Rand
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	s := &Storage{
		db: db,
		r:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	err = s.initTables()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) initTables() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS user_words (
			user_id INTEGER,
			word_id INTEGER
		);

		CREATE TABLE IF NOT EXISTS practice_questions (
			user_id INTEGER PRIMARY KEY,
			correct_id INTEGER,
			english TEXT,
			options_json TEXT
		);
	`)
	return err
}

func (s *Storage) AddWord(userID int64, wordID int) error {
	_, err := s.db.Exec(`INSERT INTO user_words (user_id, word_id) VALUES (?, ?)`, userID, wordID)
	return err
}

func (s *Storage) ClearList(userID int64) error {
	_, err := s.db.Exec(`DELETE FROM user_words WHERE user_id = ?`, userID)
	return err
}

func (s *Storage) GetUserWords(userID int64) ([]int, error) {
	rows, err := s.db.Query(`SELECT word_id FROM user_words WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wordIDs []int
	for rows.Next() {
		var wordID int
		if err := rows.Scan(&wordID); err != nil {
			return nil, err
		}
		wordIDs = append(wordIDs, wordID)
	}

	return wordIDs, nil
}

func (s *Storage) SaveQuestion(userID int64, q models.PracticeQuestion) error {
	optionsJSON, err := json.Marshal(q.Options)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO practice_questions (user_id, correct_id, english, options_json)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET correct_id=excluded.correct_id, english=excluded.english, options_json=excluded.options_json
	`, userID, q.CorrectID, q.English, string(optionsJSON))

	return err
}

func (s *Storage) GetQuestion(userID int64) (models.PracticeQuestion, bool) {
	row := s.db.QueryRow(`
		SELECT correct_id, english, options_json
		FROM practice_questions
		WHERE user_id = ?
	`, userID)

	var q models.PracticeQuestion
	var optionsJSON string

	err := row.Scan(&q.CorrectID, &q.English, &optionsJSON)
	if err != nil {
		return models.PracticeQuestion{}, false
	}

	err = json.Unmarshal([]byte(optionsJSON), &q.Options)
	if err != nil {
		return models.PracticeQuestion{}, false
	}

	return q, true
}

func (s *Storage) ClearQuestion(userID int64) error {
	_, err := s.db.Exec(`DELETE FROM practice_questions WHERE user_id = ?`, userID)
	return err
}

func (s *Storage) GeneratePracticeQuestion(userID int64) (models.PracticeQuestion, error) {
	userWordIDs, err := s.GetUserWords(userID)
	if err != nil || len(userWordIDs) == 0 {
		return models.PracticeQuestion{}, errors.New("no words available for practice")
	}

	correctIndex := s.r.Intn(len(userWordIDs))
	correctID := userWordIDs[correctIndex]
	correctWord := translation.WordByID[correctID]

	var distractors []models.Word
	for _, w := range translation.Words {
		if w.ID != correctID {
			distractors = append(distractors, w)
		}
	}
	s.r.Shuffle(len(distractors), func(i, j int) {
		distractors[i], distractors[j] = distractors[j], distractors[i]
	})

	options := append([]models.Word{correctWord}, distractors[:3]...)
	s.r.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	q := models.PracticeQuestion{
		English:   correctWord.English,
		CorrectID: correctID,
		Options:   options,
	}

	err = s.SaveQuestion(userID, q)
	if err != nil {
		return models.PracticeQuestion{}, err
	}

	return q, nil
}

func (s *Storage) IsCorrectAnswer(userID int64, wordID int) bool {
	q, ok := s.GetQuestion(userID)
	if !ok {
		return false
	}

	_ = s.ClearQuestion(userID)
	return wordID == q.CorrectID
}
