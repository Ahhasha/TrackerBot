package router

import (
	"context"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Handler interface {
	Handle(ctx context.Context, cmd *model.Command) (model.Result, error)
}
