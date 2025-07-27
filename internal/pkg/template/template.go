package template

import (
	"fmt"
	"monotonic/internal/translation"
)

func RandomWord(word translation.Word) string {
	return fmt.Sprintf("<b>random meh:</b>\n\n%s", WordCard(word))
}

func WordCard(word translation.Word) string {
	return fmt.Sprintf(
		"<b>%s</b> [%s]\n<i>%s</i>",
		word.Spanish, word.Transcription, word.English,
	)
}
