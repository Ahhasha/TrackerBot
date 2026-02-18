package expense

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Ahhasha/Tracker-bot/internal/contracts/expense"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type MonthHandler struct {
	service expense.Service
	log     *slog.Logger
}

func NewMonth(service expense.Service, log *slog.Logger) *MonthHandler {
	return &MonthHandler{
		service: service,
		log:     log,
	}
}

func (h *MonthHandler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	const op = "handler.expense.Month"

	log := h.log.With(slog.String("op", op), slog.String("cmd", string(cmd.Name)), slog.Int64("chat_id", cmd.ChatID), slog.Int64("tg_user_id", cmd.UserID))
	log.Info("handle /month")

	rep, err := h.service.Month(ctx, cmd.UserID)
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

	text := formatPeriodReport(periodMonth, rep)

	log.Info("report built", slog.String("period", "month"), slog.Int64("total", rep.Total), slog.Int("categories", len(rep.Categories)))

	return model.Result{
		ChatID: cmd.ChatID,
		Text:   text,
	}, nil

}
