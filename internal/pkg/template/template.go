package template

import (
	"fmt"
	"monotonic/internal/pkg/translation"
	"monotonic/internal/storage/models"
)

func WordCard(word models.Word) string {
	return fmt.Sprintf(
		"<b>%s</b>, <i>%s</i>\n- %s",
		word.Spanish, translation.DescribePoS(word.PoS), word.English,
	)
}
