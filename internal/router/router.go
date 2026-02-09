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
	const op = "router.Route"
	if cmd == nil {
		r.lgr.Warn("nil command", slog.String("op", op))
		return model.Result{
			ChatID: 0,
			Text:   "Unknown command",
		}
	}
	log := r.lgr.With(slog.String("op", op), slog.String("cmd", string(cmd.Name)), slog.Int64("chat_id", cmd.ChatID), slog.Int64("tg_user_id", cmd.UserID))
	h, ok := r.handlers[cmd.Name]
	if !ok {
		log.Info("unknown command")
		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "Неизвестная команда, используйте /help",
		}
	}

	res, err := h.Handle(ctx, cmd)
	if err != nil {
		log.Error("handler failed", slog.Any("err", err))
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
