package telegram

import (
	"strings"

	"github.com/Ahhasha/Tracker-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ParseUpdateToCommand(update tgbotapi.Update) *model.Command {
	if update.Message == nil {
		return nil
	}

	text := strings.TrimSpace(update.Message.Text)
	if !strings.HasPrefix(text, "/") {
		return nil
	}

	cmdText := strings.TrimPrefix(text, "/")

	var name model.CommandName

	switch cmdText {
	case "start":
		name = model.CommandStart
	case "help":
		name = model.CommandHelp
	default:
		name = model.CommandUnknown
	}

	return &model.Command{
		Name:            name,
		ChatID:          update.Message.Chat.ID,
		UserID:          update.Message.From.ID,
		UserDisplayName: update.Message.From.UserName,
		RawArgs:         "",
	}
}
