package start

import (
	"context"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
)

type RegistrationRepo interface {
	UpsertUser(ctx context.Context, db contracts.DBTX, tgID int64, username string) (userID int64, created bool, err error)
	CreateDefaultCategories(ctx context.Context, db contracts.DBTX, userID int64) (createdNames []string, err error)
}

type RegisterResult struct {
	UserID            int64
	Created           bool
	CategoriesCreated []string
}

type RegistrationService interface {
	Register(ctx context.Context, tgID int64, username string) (RegisterResult, error)
}
