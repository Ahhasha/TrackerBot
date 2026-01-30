package app

import (
	"context"
	"sync"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Ahhasha/Tracker-bot/internal/delivery/telegram"
	"github.com/Ahhasha/Tracker-bot/internal/model"
)

func (a *App) workerPool(ctx context.Context, workers int, updatesCh <-chan tgbot.Update, resultCh chan<- model.Result) {
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case update, ok := <-updatesCh:
					if !ok {
						return
					}
					cmd := telegram.ParseUpdateToCommand(update)
					if cmd == nil {
						continue
					}
					res := a.router.Route(ctx, cmd)
					select {
					case resultCh <- res:
					case <-ctx.Done():
						return
					}
				}
			}
		}(i + 1)
	}
	wg.Wait()
}
