package app

import (
	"context"
	"log/slog"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const telegramUpdateTimeoutSec = 60

func (a *App) updateReader(ctx context.Context, updatesCh chan<- tgbot.Update) error {
	const op = "app.updateReader"
	u := tgbot.NewUpdate(0)
	u.Timeout = telegramUpdateTimeoutSec
	updates := a.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return nil
		case upd, ok := <-updates:
			if !ok {
				return nil
			}
			select {
			case updatesCh <- upd:
			case <-ctx.Done():
				return nil
			default:
				a.log.Warn("updates channel full", slog.String("op", op), slog.Int("len", len(updatesCh)), slog.Int("cap", cap(updatesCh)))
			}
		}
	}
}
