package core

import (
	"context"
	"fmt"
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

func (s *spellCommand) Process(ctx context.Context, params []string) error {
	if len(params) != spellParamsLen {
		return InvalidInput
	}

	result, err := s.action(params[0])
	if err != nil {
		return err
	}

	s.show(result)
	return nil
}

func (s *spellCommand) show(letters []rune) {
	for _, letter := range letters {
		fmt.Printf("%c ", letter)
	}
	fmt.Println()
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
