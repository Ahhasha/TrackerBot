package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/delivery/telegram"
	"github.com/Ahhasha/Tracker-bot/internal/handler/start"
	"github.com/Ahhasha/Tracker-bot/internal/model"
	"github.com/Ahhasha/Tracker-bot/internal/router"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN is not set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("create bot api: %v", err)
	}

	startHandler := start.New()
	r := router.New(map[model.CommandName]router.Handler{
		model.CommandStart: startHandler,
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runBot(ctx, bot, r)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("shutting down: %v", sig)
	bot.StopReceivingUpdates()
	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	select {
	case <-done:
		log.Println("bot stopped")
	case <-shutdownCtx.Done():
		log.Println("exiting(timeout)")
	}
}

func runBot(ctx context.Context, bot *tgbotapi.BotAPI, r *router.Router) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
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
}
