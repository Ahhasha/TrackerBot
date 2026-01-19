package start

import (
	"context"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	_ = ctx

	return model.Result{
		ChatID: cmd.ChatID,
		Text:   "Приветик, курочки!",
	}, nil
}
