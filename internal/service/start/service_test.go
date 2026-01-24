package start

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/Ahhasha/Tracker-bot/internal/contracts"
)

type fakeTxManager struct {
	calls int
	err   error
}

func (f *fakeTxManager) Do(ctx context.Context, fn func(db contracts.DBTX) error) error {
	f.calls++

	if f.err != nil {
		return f.err
	}
	return fn(nil)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
}

func TestStub(t *testing.T) {}
