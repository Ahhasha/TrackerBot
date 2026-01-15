package main

import (
	"log"
	"os"

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
		log.Fatal("create bot api: %v", err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("update_id=%d chat_id=%d user_id=%d text=%q",
			update.UpdateID,
			update.Message.Chat.ID,
			update.Message.From.ID,
			update.Message.Text,
		)
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "hello")
		if _, err := bot.Send(reply); err != nil {
			log.Printf("send error: %v", err)
		}
	}
}
