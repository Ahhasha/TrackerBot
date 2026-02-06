package app

import (
	"context"

	"github.com/Ahhasha/Tracker-bot/internal/model"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *App) sender(ctx context.Context, resultsCh <-chan model.Result) {
	for {
		select {
		case <-ctx.Done():
			return
		case res, ok := <-resultsCh:
			if !ok {
				return
			}
			msg := tgbot.NewMessage(res.ChatID, res.Text)
			if _, err := a.bot.Send(msg); err != nil {
				a.log.Error("send error", "err", err, "chat_id", res.ChatID)
			}
		}
	}
}
