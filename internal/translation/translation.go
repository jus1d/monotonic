package translation

import (
	"math/rand"
	"time"
)

type Word struct {
	Spanish       string
	English       string
	Transcription string
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandom() Word {
	words := []Word{
		{"perro", "dog", "ˈpero"},
		{"gato", "cat", "ˈɡato"},
		{"libro", "book", "ˈliβɾo"},
		{"escuela", "school", "esˈkwela"},
		{"cielo", "sky", "ˈθjelo"},
		{"amar", "to love", "aˈmaɾ"},
		{"rojo", "red", "ˈroxo"},
		{"comida", "food", "koˈmiða"},
	}

	n := rng.Intn(len(words))

	return words[n]
}
