package translation

import (
	"errors"
	"math/rand"
	"monotonic/internal/pkg/models"
	"monotonic/internal/pkg/storage"

	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var wordByID = make(map[int]models.Word)

func init() {
	for _, w := range words {
		wordByID[w.ID] = w
	}
}

func GetRandomWord() models.Word {
	return words[r.Intn(len(words))]
}

func GetWordByID(id int) models.Word {
	return wordByID[id]
}

func DescribePoS(pos string) string {
	switch pos {
	case "prep":
		return "preposition"
	case "adv":
		return "adverb"
	case "adj":
		return "adjective"
	case "v":
		return "verb"
	case "nm":
		return "noun (masculine)"
	case "nf":
		return "noun (feminine)"
	case "nm/f":
		return "noun (masculine/feminine)"
	default:
		return "unknown"
	}
}

func GeneratePracticeQuestion(userID int64) (models.PracticeQuestion, error) {
	userWordIDs, ok := storage.GetUserWords(userID)
	if !ok || len(userWordIDs) == 0 {
		return models.PracticeQuestion{}, errors.New("no words available for practice")
	}

	correctIndex := r.Intn(len(userWordIDs))
	correctID := userWordIDs[correctIndex]
	correctWord := wordByID[correctID]

	var distractors []models.Word
	for _, w := range words {
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

	storage.SaveQuestion(userID, q)

	return q, nil
}

func IsCorrectAnswer(userID int64, wordID int) bool {
	q, ok := storage.GetQuestion(userID)
	if !ok {
		return false
	}

	storage.ClearQuestion(userID)
	return wordID == q.CorrectID
}
