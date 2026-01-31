package start

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	startcont "github.com/Ahhasha/Tracker-bot/internal/contracts/start"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type regService interface {
	Register(ctx context.Context, tgID int64, username string) (startcont.RegisterResult, error)
}

type Handler struct {
	lgr *slog.Logger
	reg regService
}

func New(lgr *slog.Logger, reg regService) *Handler {
	return &Handler{
		lgr: lgr,
		reg: reg,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	h.lgr.Info("start command", slog.Int64("chat_id", cmd.ChatID), slog.Int64("user_id", cmd.UserID), slog.String("username", cmd.UserDisplayName))

	res, err := h.reg.Register(ctx, cmd.UserID, cmd.UserDisplayName)
	if err != nil {
		h.lgr.Error("register user failed", slog.Any("err", err), slog.Int64("chat_id", cmd.ChatID), slog.Int64("user_id", cmd.UserID))

		return model.Result{
			ChatID: cmd.ChatID,
			Text:   "Произошла ошибка во время регистрации, пожалуйста, попробуйте чуть позже.",
		}, nil
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
