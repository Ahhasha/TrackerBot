package help

import (
	"context"

	"log/slog"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type HelpHandler struct {
	log *slog.Logger
}

func New(log *slog.Logger) *HelpHandler {
	return &HelpHandler{
		log: log,
	}
}

func (h *HelpHandler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	const op = "handler.help.Handle"

	logger := h.log.With(slog.String("op", op), slog.Int64("chat_id", cmd.ChatID))

	logger.Info("handling help command")

	text := `📊 Доступные команды:
			/start — регистрация и создание категорий
			/add <категория> <сумма> — добавить расход
			/today — расходы за сегодня
			/week — расходы за неделю
			/month — расходы за месяц
			/help — список команд`

	return model.Result{
		ChatID: cmd.ChatID,
		Text:   text,
	}, nil
}
