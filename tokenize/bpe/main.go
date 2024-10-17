package bpe

import (
	"strings"
)

func Tokenize(text string) ([]string, map[string]int) {

	text = strings.ToLower(text)

	words := strings.Fields(text)
	var (
		result      []string
		frequency   = make(map[string]int)
		maxLendWord = 15
	)

	for i := 2; i <= maxLendWord; i++ {

		for _, word := range words {
			for j := 0; j < len(word); j++ {
				if j+i > len(word) {
					continue
				}

				subWord := word[j : j+i]
				result = append(result, subWord)

				if v, ok := frequency[subWord]; ok {
					frequency[subWord] = v + 1
				} else {
					frequency[subWord] = 1
				}

			}
		}
	}

	return result, frequency

}
