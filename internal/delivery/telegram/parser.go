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

	cmdWord, args, _ := strings.Cut(cmdText, " ")
	var name model.CommandName
	switch cmdWord {
	case "start":
		name = model.CommandStart
	case "help":
		name = model.CommandHelp
	case "add":
		name = model.CommandAdd
	case "today":
		name = model.CommandToday
	case "week":
		name = model.CommandWeek
	case "month":
		name = model.CommandMonth
	default:
		name = model.CommandUnknown
	}

	return &model.Command{
		Name:            name,
		ChatID:          update.Message.Chat.ID,
		UserID:          update.Message.From.ID,
		UserDisplayName: update.Message.From.UserName,
		RawArgs:         args,
	}
}
