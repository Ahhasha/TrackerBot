package start

import (
	"context"
	"log/slog"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Handler struct {
	lgr *slog.Logger
}

func New(lgr *slog.Logger) *Handler {
	return &Handler{lgr: lgr}
}

func (h *Handler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	h.lgr.Info("start command", slog.Int64("chat_id", cmd.ChatID), slog.Int64("user_id", cmd.UserID), slog.String("username", cmd.UserDisplayName))

	return model.Result{
		ChatID: cmd.ChatID,
		Text:   "🐤🐤🐤Приветик, курочки!",
	}, nil
}
