package models

type Word struct {
	ID      int
	Spanish string
	English string
	PoS     string
}

type PracticeQuestion struct {
	English   string
	CorrectID int
	Options   []Word
}
