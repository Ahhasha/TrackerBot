package start

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type regService interface {
	RegIfNotExist(ctx context.Context, tgID int64, username string) error
}

type Handler struct {
	lgr *slog.Logger
	reg regService
}

func New(lgr *slog.Logger, reg regService) *Handler {
	return &Handler{
		lgr: lgr,
		reg: reg,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	h.lgr.Info("start command", slog.Int64("chat_id", cmd.ChatID), slog.Int64("user_id", cmd.UserID), slog.String("username", cmd.UserDisplayName))

	if err := h.reg.RegIfNotExist(ctx, cmd.UserID, cmd.UserDisplayName); err != nil {
		h.lgr.Error("register user failed", slog.Any("err", err), slog.Int64("chat_id", cmd.ChatID), slog.Int64("user_id", cmd.UserID))

		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "Произошла ошибка во время регистрации, пожалуйста, попробуйте чуть позже.",
		}, nil
	}

	return model.Result{
		ChatID: cmd.ChatID,
		Text:   fmt.Sprintf("🐤🐤🐤Приветик курочка по имени %s, ты прошёл(ла) регистрацию!", cmd.UserDisplayName),
	}, nil
}
