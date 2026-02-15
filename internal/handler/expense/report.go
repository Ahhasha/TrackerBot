package expense

import (
	"fmt"
	"strings"

	"github.com/Ahhasha/Tracker-bot/internal/model"
)

type period string

const (
	periodToday period = "today"
	periodWeek  period = "week"
	periodMonth period = "month"
)

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
