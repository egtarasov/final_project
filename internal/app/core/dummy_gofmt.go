package core

import (
	"context"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	gofmtParamsLen = 1
)

type dummyGoFmt struct {
	response string
}

func NewDummyGoFmt() *dummyGoFmt {
	return &dummyGoFmt{}
}

// Process take one argument - absolute path to .txt file.
// Add '.' at the end of every sentence and tab at the beginning of the new paragraph.
// Processing only text, which contains of latin letters.
func (d *dummyGoFmt) Process(ctx context.Context, params []string) (string, error) {
	if len(params) != gofmtParamsLen {
		return "", InvalidInput
	}

	text, err := os.ReadFile(params[0])
	if err != nil {
		return "", InvalidInput
	}

	return d.format(string(text)), nil
}

// format Немного поясню, как рабоатет форматирование:
// Будем предполагать, что текст похож на текст)
// Поэтому будем:
//  1. Избавляться от идущих подярд пробелов
//  2. Добавлять точку перед всеми заглавными буквами, независимо от того, что идет до нее(за небольшим исключением)
//  3. Добавлять tab после каждого переноса строки(новый абзац)
func (d *dummyGoFmt) format(text string) string {
	text = strings.Replace(text, "\r\n", "\n", -1)
	response := make([]rune, 0, utf8.RuneCountInString(text))
	// First paragraph tabulation
	response = append(response, '\t')
	previousRune := 'a'
	for _, r := range text {
		//Do not repeat spacebars
		if r == ' ' && previousRune == ' ' {
			response = response[:len(response)-1]
		}
		//If upper letter, add dot
		if unicode.IsUpper(r) {
			// If it first Capital letter of paragraph then just add this letter
			if len(response) == 1 || response[len(response)-1] == '\t' {
				response = append(response, r)
				continue
			}
			// base punctuational char
			punctuationSymbol := '.'
			// delete possible space symbol
			if response[len(response)-1] == ' ' {
				response = response[:len(response)-1]
			}
			// delete and recognize punctuational char
			if len(response) > 0 && unicode.IsPunct(response[len(response)-1]) {
				punctuationSymbol = response[len(response)-1]
				response = response[:len(response)-1]
			}
			// Add punctuational char and space bar
			response = append(response, punctuationSymbol, ' ')
		}

		if r == '\n' {
			// punctuational char + space - two previous chars
			noPunctuationCharacter := !(previousRune == ' ' && len(response) > 1 && unicode.IsPunct(response[len(response)-2]))
			// If punctuation char didn't appear - add dot('.').
			if !unicode.IsPunct(previousRune) && noPunctuationCharacter {
				response = append(response, '.')
			}
			// New paragraph and tab
			response = append(response, '\n', '\t')
		} else {
			response = append(response, r)
		}
		previousRune = r
	}

	return string(response)
}
