package template

import (
	"fmt"
	"monotonic/internal/pkg/models"
	"monotonic/internal/pkg/translation"
)

func WordCard(word models.Word) string {
	return fmt.Sprintf(
		"<b>%s</b>, <i>%s</i>\n- %s",
		word.Spanish, translation.DescribePoS(word.PoS), word.English,
	)
}
