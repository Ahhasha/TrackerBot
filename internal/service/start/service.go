package start

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
	startcont "github.com/Ahhasha/Tracker-bot/internal/contracts/start"
)

type service struct {
	tx   contracts.TxManager
	repo startcont.RegistrationRepo
	log  *slog.Logger
}

func NewService(tx contracts.TxManager, repo startcont.RegistrationRepo, log *slog.Logger) startcont.RegistrationService {
	return &service{
		tx:   tx,
		repo: repo,
		log:  log,
	}
}

func (s *service) Register(ctx context.Context, tgID int64, username string) (startcont.RegisterResult, error) {
	const op = "service.start.Register"
	var res startcont.RegisterResult

	err := s.tx.Do(ctx, func(db contracts.DBTX) error {
		userID, created, err := s.repo.UpsertUser(ctx, db, tgID, username)
		if err != nil {
			return fmt.Errorf("%s: upset user: %w", op, err)
		}

		res.UserID = userID
		res.Created = created

		if !created {
			return nil
		}

		cats, err := s.repo.CreateDefaultCategories(ctx, db, userID)
		if err != nil {
			return fmt.Errorf("%s: create default categories: %w", op, err)
		}

		res.CategoriesCreated = cats
		return nil
	})
	if err != nil {
		return startcont.RegisterResult{}, err
	}

	if res.Created {
		s.log.Info("user registered", slog.Int64("tg_id", tgID), slog.Int64("user_id", res.UserID))
	} else {
		s.log.Info("user already exists", slog.Int64("tg_id", tgID), slog.Int64("user_id", res.UserID))
	}
	return res, nil
}
