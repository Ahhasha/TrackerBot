package router

import "github.com/Ahhasha/Tracker-bot/internal/model"

type Router struct{}

func New() *Router {
	return &Router{}
}

func (r *Router) Route(cmd *model.Command) model.Result {
	switch cmd.Name {
	case model.CommandStart:
		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "/start ready",
		}
	case model.CommandHelp:
		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "/help find not ready",
		}
	default:
		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "Чё прислал - непонятно",
		}
	}
}
