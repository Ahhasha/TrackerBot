package main

import (
	"context"
	"log"
	"os"

	"github.com/Ahhasha/Tracker-bot/internal/delivery/telegram"
	"github.com/Ahhasha/Tracker-bot/internal/handler/start"
	"github.com/Ahhasha/Tracker-bot/internal/model"
	"github.com/Ahhasha/Tracker-bot/internal/router"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	ctx := context.Background()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("create bot api: %v", err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	startHandler := start.New()

	r := router.New(map[model.CommandName]router.Handler{
		model.CommandStart: startHandler,
	})

	for update := range updates {
		cmd := telegram.ParseUpdateToCommand(update)
		if cmd == nil {
			continue
		}

		res := r.Route(ctx, cmd)

		msg := tgbotapi.NewMessage(res.ChatID, res.Text)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("send error: %v", err)
		}
	}
}
