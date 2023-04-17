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

type DummyGoFmt struct {
	response []rune
}

func NewDummyGoFmt() *DummyGoFmt {
	return &DummyGoFmt{}
}

// Process take one argument - absolute path to .txt file.
// Add '.' at the end of every sentence and tab at the beginning of the new paragraph.
// Processing only text, which contains of latin letters.
func (d *DummyGoFmt) Process(_ context.Context, params []string) (string, error) {
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
func (d *DummyGoFmt) format(text string) string {
	text = strings.Replace(text, "\r\n", "\n", -1)
	d.response = make([]rune, 0, utf8.RuneCountInString(text))
	// First paragraph tabulation
	d.response = append(d.response, '\t')
	previousRune := 'a'
	for _, r := range text {
		//Do not repeat spacebars
		if r == ' ' && previousRune == ' ' {
			d.response = d.response[:len(d.response)-1]
		}

		//If upper letter, add dot
		if unicode.IsUpper(r) {
			d.upperLetter(previousRune)
		}

		if r == '\n' {
			d.newLine(previousRune)
		} else {
			//default action
			d.response = append(d.response, r)
		}

		previousRune = r
	}
	return string(d.response)
}

func (d *DummyGoFmt) upperLetter(previousRune rune) {
	// Do nothing if previous letter is upper
	if unicode.IsUpper(previousRune) {
		return
	}
	// If it first Capital letter of paragraph then just add this letter
	if len(d.response) == 1 || d.response[len(d.response)-1] == '\t' {
		return
	}
	// base punctuational char
	punctuationSymbol := '.'
	// delete possible space symbol
	if d.response[len(d.response)-1] == ' ' {
		d.response = d.response[:len(d.response)-1]
	}
	// delete and recognize punctuational char
	if len(d.response) > 0 && unicode.IsPunct(d.response[len(d.response)-1]) {
		punctuationSymbol = d.response[len(d.response)-1]
		d.response = d.response[:len(d.response)-1]
	}
	// Add punctuational char and space bar
	d.response = append(d.response, punctuationSymbol, ' ')
}

func (d *DummyGoFmt) newLine(previousRune rune) {
	// If this is just empty line - skip
	if previousRune == '\n' {
		d.response = append(d.response, '\n', '\t')
	} else {
		// punctuational char + space - two previous chars
		noPunctuationCharacter := !(previousRune == ' ' && len(d.response) > 1 && unicode.IsPunct(d.response[len(d.response)-2]))
		// If punctuation char didn't appear - add dot('.').
		if !unicode.IsPunct(previousRune) && noPunctuationCharacter {
			d.response = append(d.response, '.')
		}
		// New paragraph and tab
		d.response = append(d.response, '\n', '\t')
	}
}
