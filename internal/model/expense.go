package model

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidAmount   = errors.New("сумма должна быть положительной")
	ErrInvalidCategory = errors.New("категория обязательна")
	ErrInvalidUser     = errors.New("пользователь обязателен")
)

type Expense struct {
	ID          int64
	UserID      int64
	Amount      int64
	CategoryID  int64
	Description string
	CreatedAt   time.Time
}

func (e *Expense) Validate() error {
	if e.Amount < +0 {
		return ErrInvalidAmount
	}

	if e.CategoryID <= 0 {
		return ErrInvalidCategory
	}

	if e.UserID <= 0 {
		return ErrInvalidUser
	}

	return nil
}

func (e *Expense) Format() string {
	return fmt.Sprintf("%d₽", e.Amount)
}
