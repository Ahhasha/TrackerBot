package expense

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPeriodRange_Today(t *testing.T) {
	now := time.Date(2026, 2, 13, 10, 30, 0, 0, time.UTC)

	start, end := periodRange(periodToday, now)

	require.Equal(t, time.Date(2026, 2, 13, 0, 0, 0, 0, time.UTC), start)
	require.Equal(t, time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC), end)
}

func TestPeriodRange_Week_Monday(t *testing.T) {
	now := time.Date(2026, 2, 9, 15, 0, 0, 0, time.UTC)
	start, end := periodRange(periodWeek, now)

	require.Equal(t, time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC), start)
	require.Equal(t, time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC), end)
}

func TestPeriodRange_Week_Sunday(t *testing.T) {
	now := time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC)
	start, end := periodRange(periodWeek, now)

	require.Equal(t, time.Date(2026, 2, 9, 0, 0, 0, 0, time.UTC), start)
	require.Equal(t, time.Date(2026, 2, 16, 0, 0, 0, 0, time.UTC), end)
}

func TestPeriodRange_Month(t *testing.T) {
	now := time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	start, end := periodRange(periodMonth, now)

	require.Equal(t, time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC), start)
	require.Equal(t, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), end)
}
