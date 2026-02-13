package expense

import (
	"testing"
	"time"

	"github.com/Ahhasha/Tracker-bot/internal/model"
	"github.com/stretchr/testify/require"
)

func TestBuildPeriodReport(t *testing.T) {
	anchor := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)

	rows := []model.ExpenseWithCategory{
		{Amount: 300, Description: "метро", Category: "Транспорт"},
		{Amount: 500, Description: "кофе", Category: "Еда"},
		{Amount: 700, Description: "обед", Category: "Еда"},
		{Amount: 500, Description: "такси", Category: "Транспорт"},
	}

	got := buildPeriodReport(anchor, rows)

	require.Equal(t, anchor, got.Date)

	require.Equal(t, int64(2000), got.Total)

	require.Len(t, got.Categories, 2)

	require.Equal(t, "Еда", got.Categories[0].Name)
	require.Equal(t, int64(1200), got.Categories[0].Total)

	require.Equal(t, "Транспорт", got.Categories[1].Name)
	require.Equal(t, int64(800), got.Categories[1].Total)

	require.Len(t, got.Categories[0].Items, 2)
	require.Equal(t, int64(500), got.Categories[0].Items[0].Amount)
	require.Equal(t, "кофе", got.Categories[0].Items[0].Description)
}

func TestBuildPeriodReport_SortByName(t *testing.T) {
	anchor := time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC)

	rows := []model.ExpenseWithCategory{
		{Amount: 100, Description: "x", Category: "Б"},
		{Amount: 100, Description: "y", Category: "А"},
	}

	got := buildPeriodReport(anchor, rows)

	require.Len(t, got.Categories, 2)
	require.Equal(t, "А", got.Categories[0].Name)
	require.Equal(t, "Б", got.Categories[1].Name)
}
