package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/app"
	"github.com/Ahhasha/Tracker-bot/internal/config"
	"github.com/Ahhasha/Tracker-bot/internal/database"
	expenseHandler "github.com/Ahhasha/Tracker-bot/internal/handler/expense"
	"github.com/Ahhasha/Tracker-bot/internal/handler/start"
	"github.com/Ahhasha/Tracker-bot/internal/model"
	categoryRepo "github.com/Ahhasha/Tracker-bot/internal/repository/category"
	expenseRepo "github.com/Ahhasha/Tracker-bot/internal/repository/expense"
	repoStart "github.com/Ahhasha/Tracker-bot/internal/repository/start"
	userRepo "github.com/Ahhasha/Tracker-bot/internal/repository/user"
	"github.com/Ahhasha/Tracker-bot/internal/router"
	expenseService "github.com/Ahhasha/Tracker-bot/internal/service/expense"
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

	txManager := database.NewTxManager(pool)

	startRepository := repoStart.NewPostgresRepo()
	userRepository := userRepo.NewPostgresRepo()
	expenseRepository := expenseRepo.NewPostgresRepo()
	categoryRepository := categoryRepo.NewPostgresRepo()

	regService := serv.NewService(txManager, startRepository, logger)
	expenseService := expenseService.NewService(txManager, expenseRepository, categoryRepository, userRepository, time.Now)

	startHandler := start.New(logger, regService)
	addHandler := expenseHandler.NewAdd(expenseService, logger)
	todayHandler := expenseHandler.NewToday(expenseService, logger)
	weekHandler := expenseHandler.NewWeek(expenseService, logger)
	monthHandler := expenseHandler.NewMonth(expenseService, logger)

	r := router.New(map[model.CommandName]router.Handler{
		model.CommandStart: startHandler,
		model.CommandAdd:   addHandler,
		model.CommandToday: todayHandler,
		model.CommandWeek:  weekHandler,
		model.CommandMonth: monthHandler,
	}, logger)

	app := app.NewBot(bot, r, logger)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Run(ctx); err != nil {
			logger.Error("app stopped with error", slog.Any("err", err))
			cancel()
		}
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
