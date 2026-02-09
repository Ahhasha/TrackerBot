package expense

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Ahhasha/Tracker-bot/internal/contracts/expense"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type WeekHandler struct {
	service expense.Service
	logger  *slog.Logger
}

func NewWeek(service expense.Service, logger *slog.Logger) *WeekHandler {
	return &WeekHandler{
		service: service,
		logger:  logger,
	}
}

func (h *WeekHandler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	const op = "handler.expense.Week"

	log := h.logger.With(slog.String("op", op), slog.String("cmd", string(cmd.Name)), slog.Int64("chat_id", cmd.ChatID), slog.Int64("tg_user_id", cmd.UserID))
	log.Info("handle /week")

	rep, err := h.service.Week(ctx, cmd.UserID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotRegistered):
			return model.Result{
				ChatID: cmd.ChatID,
				Text:   "❌ Вы не зарегистрированы. Используйте /start",
			}, nil
		}
		return model.Result{}, err
	}

	text := formatPeriodReport(periodWeek, rep)

	log.Info("report built", slog.String("period", "week"), slog.Int64("total", rep.Total), slog.Int("categories", len(rep.Categories)))

	return model.Result{
		ChatID: cmd.ChatID,
		Text:   text,
	}, nil
}
