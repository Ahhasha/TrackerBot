package expense

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Ahhasha/Tracker-bot/internal/contracts/expense"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type TodayHandler struct {
	service expense.Service
	logger  *slog.Logger
}

func NewToday(service expense.Service, logger *slog.Logger) *TodayHandler {
	return &TodayHandler{
		service: service,
		logger:  logger,
	}
}

type period string

const (
	periodToday period = "today"
	periodWeek  period = "week"
	periodMonth period = "month"
)

func (h *TodayHandler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	const op = "handler.expense.Today"

	log := h.logger.With(slog.String("op", op), slog.String("cmd", string(cmd.Name)), slog.Int64("chat_id", cmd.ChatID), slog.Int64("tg_user_id", cmd.UserID))
	log.Info("handle /today")

	rep, err := h.service.Today(ctx, cmd.UserID)
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

	text := formatPeriodReport(periodToday, rep)

	log.Info("report built", slog.String("period", "today"), slog.Int64("total", rep.Total), slog.Int("categories", len(rep.Categories)))

	return model.Result{
		ChatID: cmd.ChatID,
		Text:   text,
	}, nil
}
