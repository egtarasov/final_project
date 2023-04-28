package core

import (
	"context"
	"strings"
	"unicode"
)

const (
	spellParamsLen = 1
)

type spellCommand struct {
}

func NewSpellCommand() *spellCommand {
	return &spellCommand{}
}

func (s *spellCommand) Process(ctx context.Context, params []string) (string, error) {
	if len(params) != spellParamsLen {
		return "", InvalidInput
	}

	result, err := s.action(params[0])
	if err != nil {
		return "", err
	}

	return s.show(result), nil
}

func (s *spellCommand) show(letters []rune) string {
	var stringBuilder strings.Builder
	for _, letter := range letters {
		stringBuilder.WriteRune(letter)
		stringBuilder.WriteRune(' ')
	}
	stringBuilder.WriteRune('\n')
	return stringBuilder.String()
}

func (s *spellCommand) action(word string) ([]rune, error) {
	result := make([]rune, len(word))
	for _, letter := range word {
		if !unicode.IsLetter(letter) {
			return nil, InvalidInput
		}
		result = append(result, letter)
	}

	return result, nil
}
