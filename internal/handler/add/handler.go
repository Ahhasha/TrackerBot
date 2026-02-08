package add

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Ahhasha/Tracker-bot/internal/contracts/add"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Handler struct {
	service add.AddService
	logger  *slog.Logger
}

func New(service add.AddService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	h.logger.Info("start command", slog.Int64("chat_id", cmd.ChatID), slog.Int64("user_id", cmd.UserID), slog.String("username", cmd.UserDisplayName))

	const op = "handler.add.Handle"
	req, err := h.parseAddArgs(cmd.RawArgs)
	if err != nil {
		h.logger.Warn("fail parse /add arguments", slog.String("error", err.Error()), slog.String("args", cmd.RawArgs))

		return model.Result{
			ChatID: cmd.ChatID,
			Text:   h.formatUsage(),
		}, nil
	}

	expenseID, err := h.service.AddExpense(ctx, cmd.UserID, req)
	if err != nil {
		h.logger.Error("fail add expense", slog.String("op", op), slog.Int64("tguser_id", cmd.UserID), slog.Any("error", err))

		errorText := h.error(err, req.Category)
		return model.Result{
			ChatID: cmd.ChatID,
			Text:   errorText,
		}, nil
	}

	successText := fmt.Sprintf("✅ Расход добавлен!\n\n"+
		"💰 Сумма: %d ₽\n"+
		"📂 Категория: %s\n"+
		"📝 ID: #%d",
		req.Amount, req.Category, expenseID)

	if req.Description != "" {
		successText += fmt.Sprintf("\n📄 Описание: %s", req.Description)
	}

	h.logger.Info("expense added successfully", slog.Int64("expense_id", expenseID), slog.Int64("tguser_id", cmd.UserID), slog.Int64("amount", req.Amount))

	return model.Result{
		ChatID: cmd.ChatID,
		Text:   successText,
	}, nil
}

func (h *Handler) parseAddArgs(rawArgs string) (add.AddRequest, error) {
	parts := strings.Fields(rawArgs)

	if len(parts) < 2 {
		return add.AddRequest{}, fmt.Errorf("недостаточно аргументов")
	}

	amount, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return add.AddRequest{}, fmt.Errorf("неверный формат суммы: %v", err)
	}

	if amount <= 0 {
		return add.AddRequest{}, fmt.Errorf("сумма должна быть положительной")
	}

	category := parts[1]

	var description string
	if len(parts) > 2 {
		description = strings.Join(parts[2:], " ")
	}

	return add.AddRequest{
		Amount:      amount,
		Category:    category,
		Description: description,
	}, nil
}

func (h *Handler) formatUsage() string {
	return `❌ Неверный формат команды.

Используйте: /add <сумма> <категория> [описание]

Примеры:
/add 1500 еда обед в кафе
/add 500 транспорт такси до работы
/add 3000 развлечения кино`

}

func (h *Handler) error(err error, category string) string {
	errorMsg := err.Error()

	if strings.Contains(errorMsg, "category not found") || strings.Contains(errorMsg, "категория не найдена") {
		return fmt.Sprintf("❌ Категория '%s' не найдена.\n\n"+
			"Используйте /categories чтобы увидеть доступные категории.", category)
	}

	if strings.Contains(errorMsg, "сумма должна быть положительной") {
		return "❌ Сумма должна быть положительным числом."
	}

	if strings.Contains(errorMsg, "категория обязательна") {
		return "❌ Категория обязательна. Укажите категорию расхода."
	}

	return "❌ Не удалось добавить расход. Попробуйте позже."
}
