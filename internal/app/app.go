package app

import (
	"context"
	"log/slog"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type App struct {
	bot   *tgbot.BotAPI
	model *model.Result
	log   *slog.Logger
}

func NewBot(bot *tgbot.BotAPI, m *model.Result, log *slog.Logger) *App {
	return &App{
		bot:   bot,
		model: m,
		log:   log,
	}
}

func (a *App) Run(ctx context.Context) error {
	updatesCh := make(chan tgbot.Update, 100)
	resultsCh := make(chan model.Result, 100)
	go a.sender(ctx, resultsCh)

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.updateReader(ctx, updatesCh)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
