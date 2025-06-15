package web

import (
	"context"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Key int

const (
	Localizer Key = iota
	TraceID
)

func localizerFrom(ctx context.Context) (*i18n.Localizer, error) {
	loc, ok := ctx.Value(Localizer).(*i18n.Localizer)
	if !ok {
		return nil, fmt.Errorf("could not get localizer from context")
	}
	return loc, nil
}

