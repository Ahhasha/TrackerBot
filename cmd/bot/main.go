package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/config"
	"github.com/Ahhasha/Tracker-bot/internal/database"
	"github.com/Ahhasha/Tracker-bot/internal/delivery/telegram"
	"github.com/Ahhasha/Tracker-bot/internal/handler/start"
	"github.com/Ahhasha/Tracker-bot/internal/model"
	repoStart "github.com/Ahhasha/Tracker-bot/internal/repository/start"
	"github.com/Ahhasha/Tracker-bot/internal/router"
	serv "github.com/Ahhasha/Tracker-bot/internal/service/start"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if err := godotenv.Load(); err != nil {
		logger.Info(".env file not found")
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid config", slog.Any("err", err))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		logger.Error("create bot api:", slog.Any("err", err))
		os.Exit(1)
	}

	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		logger.Error("db init failed", slog.Any("err", err))
		os.Exit(1)
	}
	defer pool.Close()

	repo := repoStart.NewRepo(pool)
	regService := serv.NewService(pool, logger, repo)
	startHandler := start.New(logger, regService)

	r := router.New(map[model.CommandName]router.Handler{
		model.CommandStart: startHandler,
	}, logger)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runBot(ctx, bot, r, logger)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("shutting down", slog.String("signal", sig.String()))
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
		logger.Info("bot stopped")
	case <-shutdownCtx.Done():
		logger.Warn("exiting(timeout)")
	}
}

func runBot(ctx context.Context, bot *tgbotapi.BotAPI, r *router.Router, logger *slog.Logger) {
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
				logger.Error("send error", slog.Any("err", err), slog.Int64("chat_id", res.ChatID))
			}
		}
	}
}
