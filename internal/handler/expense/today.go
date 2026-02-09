package expense

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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

func formatPeriodReport(p period, rep model.PeriodReport) string {
	title := periodTitle(p)

	var b strings.Builder
	b.WriteString(title)
	b.WriteString("\n\n")

	if rep.Total == 0 || len(rep.Categories) == 0 {
		b.WriteString("Расходов нет 🙂")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("Итого: %d ₽\n\n", rep.Total))

	for _, c := range rep.Categories {
		b.WriteString(fmt.Sprintf("📂 %s — %d ₽\n", c.Name, c.Total))

		for _, it := range c.Items {
			desc := strings.TrimSpace(it.Description)
			if desc != "" {
				b.WriteString(fmt.Sprintf("  • %d ₽ — %s\n", it.Amount, desc))
			} else {
				b.WriteString(fmt.Sprintf("  • %d ₽\n", it.Amount))
			}
		}

		b.WriteString("\n")
	}

	return strings.TrimRight(b.String(), "\n")
}

func periodTitle(p period) string {
	switch p {
	case periodToday:
		return "📅 Отчёт за сегодня"
	case periodWeek:
		return "📅 Отчёт за неделю"
	case periodMonth:
		return "📅 Отчёт за месяц"
	default:
		return "📅 Отчёт"
	}
}
