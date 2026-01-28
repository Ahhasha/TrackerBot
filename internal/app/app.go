package app

import (
	"context"
	"log/slog"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Ahhasha/Tracker-bot/internal/router"
)

type App struct {
	bot    *tgbot.BotAPI
	router *router.Router
	log    *slog.Logger
}

func NewBot(bot *tgbot.BotAPI, r *router.Router, log *slog.Logger) *App {
	return &App{
		bot:    bot,
		router: r,
		log:    log,
	}
}

func (a *App) Run(ctx context.Context) error {
	_ = ctx
	return nil
}
