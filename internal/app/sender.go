package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/model"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	maxSendRetries = 3
	retryDelay     = 500 * time.Millisecond
)

func (a *App) sender(ctx context.Context, resultsCh <-chan model.Result) error {
	const op = "app.sender"

	for {
		select {
		case <-ctx.Done():
			return nil
		case res, ok := <-resultsCh:
			if !ok {
				return nil
			}
			if res.ChatID == 0 {
				a.log.Warn("skip sending message: empty chat id", slog.String("op", op))
				continue
			}

			msg := tgbot.NewMessage(res.ChatID, res.Text)

			var err error

			for attempt := 1; attempt <= maxSendRetries; attempt++ {
				_, err = a.bot.Send(msg)
				if err == nil {
					break
				}

				a.log.Warn("send failed, retrying", slog.String("op", op), slog.Int("attempt", attempt), slog.Int64("chat_id", res.ChatID), slog.Any("err", err))

				select {
				case <-time.After(retryDelay):
				case <-ctx.Done():
					return nil
				}
			}

			if err != nil {
				a.log.Error("send permanently failed", slog.String("op", op), slog.Int64("chat_id", res.ChatID), slog.Any("err", err))
			}
		}
	}
}
