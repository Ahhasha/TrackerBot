package start

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	startcont "github.com/Ahhasha/Tracker-bot/internal/contracts/start"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Handler struct {
	lgr *slog.Logger
	reg startcont.RegistrationService
}

func New(lgr *slog.Logger, reg startcont.RegistrationService) *Handler {
	return &Handler{
		lgr: lgr,
		reg: reg,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	const op = "handler.start.Handle"
	log := h.lgr.With(slog.String("op", op), slog.String("cmd", string(cmd.Name)), slog.Int64("chat_id", cmd.ChatID), slog.Int64("tg_user_id", cmd.UserID))
	log.Info("start command")

	res, err := h.reg.Register(ctx, cmd.UserID, cmd.UserDisplayName)
	if err != nil {
		return model.Result{}, err
	}

	if res.Created {
		var b strings.Builder
		for _, c := range res.CategoriesCreated {
			b.WriteString("• ")
			b.WriteString(c)
			b.WriteString("\n")
		}
		catsText := strings.TrimRight(b.String(), "\n")

		return model.Result{
			ChatID: cmd.ChatID,
			Text: fmt.Sprintf("👋 Добро пожаловать в Expense Tracker!\n"+
				"Я помогу вам отслеживать расходы и управлять бюджетами.\n\n"+
				"✅ Вы зарегистрированы!\n"+
				"📂 Созданы базовые категории:\n"+
				"%s\n\n"+
				"Используйте /help для списка команд.",
				catsText,
			),
		}, nil
	}
	return model.Result{
		ChatID: cmd.ChatID,
		Text:   "✅ Вы уже зарегистрированы! Используйте /help для списка команд.",
	}, nil
}
