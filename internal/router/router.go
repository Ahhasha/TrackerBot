package router

import (
	"context"
	"log/slog"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Router struct {
	handlers map[model.CommandName]Handler
	lgr      *slog.Logger
}

func New(handlers map[model.CommandName]Handler, lgr *slog.Logger) *Router {
	return &Router{
		handlers: handlers,
		lgr:      lgr,
	}
}

func (r *Router) Route(ctx context.Context, cmd *model.Command) model.Result {
	if cmd == nil {
		return model.Result{
			ChatID: 0,
			Text:   "Unknown command",
		}
	}

	h, ok := r.handlers[cmd.Name]
	if !ok {
		r.lgr.Info("unknown command", slog.String("cmd", string(cmd.Name)), slog.Int64("chat_id", cmd.ChatID), slog.Int64("user_id", cmd.UserID))
		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "Неизвестная команда, используйте /help",
		}
	}

	res, err := h.Handle(ctx, cmd)
	if err != nil {
		r.lgr.Error("Handler error", slog.Any("err", err), slog.String("cmd", string(cmd.Name)), slog.Int64("chat_id", cmd.ChatID), slog.Int64("user_id", cmd.UserID))
		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "Произошла ошибка, попробуйте позже.",
		}
	}

	if res.ChatID == 0 {
		res.ChatID = cmd.ChatID
	}
	return res
}
