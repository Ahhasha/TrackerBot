package router

import (
	"context"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Router struct {
	handlers map[model.CommandName]Handler
}

func New(handlers map[model.CommandName]Handler) *Router {
	if handlers == nil {
		handlers = make(map[model.CommandName]Handler)
	}
	return &Router{handlers: handlers}
}

func (r *Router) Route(ctx context.Context, cmd *model.Command) model.Result {
	if cmd == nil {
		return model.Result{
			ChatID: 0,
			Text:   "Unknown command",
		}
	}

	h, ok := r.handlers[cmd.Name]
	if !ok || h == nil {
		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "Unknown command, use /help",
		}
	}

	res, err := h.Handle(ctx, cmd)
	if err != nil {
		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "Error. Please try again later.",
		}
	}

	if res.ChatID == 0 {
		res.ChatID = cmd.ChatID
	}
	return res
}
