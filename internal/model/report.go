package model

import "time"

type PeriodReport struct {
	Date       time.Time
	Categories []CategoryReport
	Total      int64
}

type CategoryReport struct {
	Name  string
	Total int64
	Items []ExpenseItem
}

type ExpenseItem struct {
	Amount      int64
	Description string
}

type ExpenseWithCategory struct {
	Amount      int64
	Description string
	Category    string
	CreatedAt   time.Time
}
