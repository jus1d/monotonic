package storage

var userWords = make(map[int][]int)

func AddWord(userID int, wordID int) {
	words, ok := userWords[userID]
	if ok {
		words = append(words, wordID)
		userWords[userID] = words
	} else {
		words = make([]int, 1)
		words = append(words, wordID)
	}
}

func GetUserWords(userID int) ([]int, bool) {
	words, ok := userWords[userID]
	return words, ok
}
