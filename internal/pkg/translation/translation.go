package translation

import (
	"math/rand"
	"time"
)

type Word struct {
	Spanish string
	English string
	// Transcription string
	PoS string
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandom() Word {
	n := rng.Intn(len(words))

	return words[n]
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
