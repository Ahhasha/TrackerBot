package expense

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Ahhasha/Tracker-bot/internal/contracts/expense"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type Handler struct {
	service expense.ExpenseService
	logger  *slog.Logger
}

func New(service expense.ExpenseService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *model.Command) (model.Result, error) {
	const op = "handler.add.Handle"
	log := h.logger.With(slog.String("op", op), slog.String("cmd", string(cmd.Name)), slog.Int64("chat_id", cmd.ChatID), slog.Int64("tg_user_id", cmd.UserID))
	log.Info("handle /add")

	req, err := h.parseAddArgs(cmd.RawArgs)
	if err != nil {
		log.Warn("fail parse /add arguments", slog.Any("err", err), slog.String("args", cmd.RawArgs))

		return model.Result{
			ChatID: cmd.ChatID,
			Text: `❌ Неверный формат команды.
				Используйте: /add <сумма> <категория> [описание]
				Примеры:
				/add 1500 еда обед в кафе
				/add 500 транспорт такси до работы
				/add 3000 развлечения кино`,
		}, nil
	}

	expenseID, err := h.service.AddExpense(ctx, cmd.UserID, req)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrUserNotRegistered):
			return model.Result{
				ChatID: cmd.ChatID,
				Text:   "❌ Вы не зарегистрированы. Используйте /start",
			}, nil

		case errors.Is(err, model.ErrCategoryNotFound):
			return model.Result{
				ChatID: cmd.ChatID,
				Text:   "❌ Категория не найдена. Используйте /categories",
			}, nil

		case errors.Is(err, model.ErrInvalidAmount):
			return model.Result{
				ChatID: cmd.ChatID,
				Text:   "❌ Ошибка: сумма должна быть положительным числом",
			}, nil

		case errors.Is(err, model.ErrInvalidCategory):
			return model.Result{
				ChatID: cmd.ChatID,
				Text:   "❌ Категория не указана. Пример: /add 500 Еда Обед",
			}, nil
		}
		return model.Result{}, err
	}

	successText := fmt.Sprintf("✅ Расход добавлен!\n\n"+
		"💰 Сумма: %d ₽\n"+
		"📂 Категория: %s\n"+
		"📝 ID: #%d",
		req.Amount, req.Category, expenseID)

	if req.Description != "" {
		successText += fmt.Sprintf("\n📄 Описание: %s", req.Description)
	}

	log.Info("expense added successfully", slog.Int64("expense_id", expenseID), slog.Int64("amount", req.Amount))

	return model.Result{
		ChatID: cmd.ChatID,
		Text:   successText,
	}, nil
}

func (h *Handler) parseAddArgs(rawArgs string) (expense.ExpenseRequest, error) {
	parts := strings.Fields(rawArgs)

	if len(parts) < 2 {
		return expense.ExpenseRequest{}, fmt.Errorf("недостаточно аргументов")
	}

	amount, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return expense.ExpenseRequest{}, fmt.Errorf("неверный формат суммы: %v", err)
	}

	category := parts[1]

	var description string
	if len(parts) > 2 {
		description = strings.Join(parts[2:], " ")
	}

	return expense.ExpenseRequest{
		Amount:      amount,
		Category:    category,
		Description: description,
	}, nil
}
