package core

import (
	"context"
	"os"
)

const (
	gofmtParamsLen = 1
)

type dummyGoFmt struct {
}

func (d *dummyGoFmt) Process(ctx context.Context, params []string) error {
	if len(params) != gofmtParamsLen {
		return InvalidInput
	}
	text, err := os.ReadFile(params[0])
	if err != nil {
		return InvalidInput
	}
	_ = text
	return nil
}
