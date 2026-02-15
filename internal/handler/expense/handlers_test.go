package expense_test

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/contracts/mocks"
	expenseHandler "github.com/Ahhasha/Tracker-bot/internal/handler/expense"
	"github.com/Ahhasha/Tracker-bot/internal/model"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestAddHandler_ParseError_ServiceNotCalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockService(ctrl)

	svc.EXPECT().AddExpense(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	h := expenseHandler.NewAdd(svc, newTestLogger())

	cmd := &model.Command{
		Name:    model.CommandAdd,
		ChatID:  100,
		UserID:  200,
		RawArgs: "нет_двух_аргументов",
	}

	got, err := h.Handle(context.Background(), cmd)
	require.NoError(t, err)

	require.Equal(t, cmd.ChatID, got.ChatID)
	require.Contains(t, got.Text, "❌ Неверный формат команды")
	require.Contains(t, got.Text, "Используйте: /add")
}

func TestAddHandler_UserNotRegistered_MappedMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockService(ctrl)
	h := expenseHandler.NewAdd(svc, newTestLogger())

	cmd := &model.Command{
		Name:    model.CommandAdd,
		ChatID:  100,
		UserID:  777,
		RawArgs: "500 еда обед",
	}

	svc.EXPECT().AddExpense(gomock.Any(), cmd.UserID, gomock.Any()).Return(int64(0), model.ErrUserNotRegistered)

	got, err := h.Handle(context.Background(), cmd)
	require.NoError(t, err)

	require.Equal(t, cmd.ChatID, got.ChatID)
	require.Equal(t, "❌ Вы не зарегистрированы. Используйте /start", got.Text)
}

func TestAddHandler_HappyPath_WithDescription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockService(ctrl)
	h := expenseHandler.NewAdd(svc, newTestLogger())

	cmd := &model.Command{
		Name:    model.CommandAdd,
		ChatID:  100,
		UserID:  777,
		RawArgs: "500 еда обед в кафе",
	}

	svc.EXPECT().AddExpense(gomock.Any(), cmd.UserID, gomock.Any()).Return(int64(42), nil)

	got, err := h.Handle(context.Background(), cmd)
	require.NoError(t, err)

	require.Equal(t, cmd.ChatID, got.ChatID)
	require.Contains(t, got.Text, "✅ Расход добавлен!")
	require.Contains(t, got.Text, "💰 Сумма: 500 ₽")
	require.Contains(t, got.Text, "📂 Категория: еда")
	require.Contains(t, got.Text, "📝 ID: #42")
	require.Contains(t, got.Text, "📄 Описание: обед в кафе")
}

func TestTodayHandler_UserNotRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockService(ctrl)
	h := expenseHandler.NewToday(svc, newTestLogger())

	cmd := &model.Command{
		Name:   model.CommandToday,
		ChatID: 10,
		UserID: 20,
	}

	svc.EXPECT().Today(gomock.Any(), cmd.UserID).Return(model.PeriodReport{}, model.ErrUserNotRegistered)

	got, err := h.Handle(context.Background(), cmd)
	require.NoError(t, err)

	require.Equal(t, cmd.ChatID, got.ChatID)
	require.Equal(t, "❌ Вы не зарегистрированы. Используйте /start", got.Text)
}

func TestTodayHandler_HappyPath_EmptyReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockService(ctrl)
	h := expenseHandler.NewToday(svc, newTestLogger())

	cmd := &model.Command{
		Name:   model.CommandToday,
		ChatID: 10,
		UserID: 20,
	}

	rep := model.PeriodReport{
		Date:       time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC),
		Categories: nil,
		Total:      0,
	}

	svc.EXPECT().Today(gomock.Any(), cmd.UserID).Return(rep, nil)

	got, err := h.Handle(context.Background(), cmd)
	require.NoError(t, err)

	require.Equal(t, cmd.ChatID, got.ChatID)

	require.Contains(t, got.Text, "📅 Отчёт за сегодня")
	require.Contains(t, got.Text, "Расходов нет 🙂")
}

func TestWeekHandler_UserNotRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockService(ctrl)
	h := expenseHandler.NewWeek(svc, newTestLogger())

	cmd := &model.Command{
		Name:   model.CommandWeek,
		ChatID: 10,
		UserID: 20,
	}

	svc.EXPECT().Week(gomock.Any(), cmd.UserID).Return(model.PeriodReport{}, model.ErrUserNotRegistered)

	got, err := h.Handle(context.Background(), cmd)
	require.NoError(t, err)

	require.Equal(t, cmd.ChatID, got.ChatID)
	require.Equal(t, "❌ Вы не зарегистрированы. Используйте /start", got.Text)
}

func TestWeekHandler_HappyPath_EmptyReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockService(ctrl)
	h := expenseHandler.NewWeek(svc, newTestLogger())

	cmd := &model.Command{
		Name:   model.CommandWeek,
		ChatID: 10,
		UserID: 20,
	}

	rep := model.PeriodReport{
		Date:  time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC),
		Total: 0,
	}

	svc.EXPECT().Week(gomock.Any(), cmd.UserID).Return(rep, nil)

	got, err := h.Handle(context.Background(), cmd)
	require.NoError(t, err)

	require.Equal(t, cmd.ChatID, got.ChatID)
	require.Contains(t, got.Text, "📅 Отчёт за неделю")
	require.Contains(t, got.Text, "Расходов нет 🙂")
}

func TestMonthHandler_UserNotRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockService(ctrl)
	h := expenseHandler.NewMonth(svc, newTestLogger())

	cmd := &model.Command{
		Name:   model.CommandMonth,
		ChatID: 10,
		UserID: 20,
	}

	svc.EXPECT().Month(gomock.Any(), cmd.UserID).Return(model.PeriodReport{}, model.ErrUserNotRegistered)

	got, err := h.Handle(context.Background(), cmd)
	require.NoError(t, err)

	require.Equal(t, cmd.ChatID, got.ChatID)
	require.Equal(t, "❌ Вы не зарегистрированы. Используйте /start", got.Text)
}

func TestMonthHandler_HappyPath_EmptyReport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockService(ctrl)
	h := expenseHandler.NewMonth(svc, newTestLogger())

	cmd := &model.Command{
		Name:   model.CommandMonth,
		ChatID: 10,
		UserID: 20,
	}

	rep := model.PeriodReport{
		Date:  time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC),
		Total: 0,
	}

	svc.EXPECT().Month(gomock.Any(), cmd.UserID).Return(rep, nil)

	got, err := h.Handle(context.Background(), cmd)
	require.NoError(t, err)

	require.Equal(t, cmd.ChatID, got.ChatID)
	require.Contains(t, got.Text, "📅 Отчёт за месяц")
	require.Contains(t, got.Text, "Расходов нет 🙂")
}
