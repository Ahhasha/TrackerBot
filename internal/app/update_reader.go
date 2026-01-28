package app

import (
	"context"
	"fmt"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *App) updateReader(ctx context.Context, updatesCh chan<- tgbot.Update) error {
	const op = "app.updateReader"
	u := tgbot.NewUpdate(0)
	u.Timeout = 60
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
				// апдейт ушёл воркерам (позже)
			case <-ctx.Done():
				return nil
			default:
				return fmt.Errorf("%s: full channel", op)
			}
		}
	}
}
