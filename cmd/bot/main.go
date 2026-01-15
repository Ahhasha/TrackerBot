package main

import (
	"log"
	"os"

	"github.com/Ahhasha/Tracker-bot/internal/delivery/telegram"
	"github.com/Ahhasha/Tracker-bot/internal/router"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("create bot api: %v", err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	r := router.New()

	for update := range updates {
		cmd := telegram.ParseUpdateToCommand(update)
		if cmd == nil {
			continue
		}

		log.Printf("command=%s chat_id=%d user_id=%d args=%q", cmd.Name, cmd.ChatID, cmd.UserID, cmd.RawArgs)

		res := r.Route(cmd)

		msg := tgbotapi.NewMessage(res.ChatID, res.Text)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("send error: %v", err)
		}
	}
}
