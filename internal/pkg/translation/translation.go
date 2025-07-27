package translation

import (
	"math/rand"
	"monotonic/internal/storage/models"

	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var WordByID = make(map[int]models.Word)

func init() {
	for _, w := range Words {
		WordByID[w.ID] = w
	}
}

func GetRandomWord() models.Word {
	return Words[r.Intn(len(Words))]
}

func GetWordByID(id int) models.Word {
	return WordByID[id]
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
	case "conj":
		return "conjunction"
	default:
		return pos
	}
}
