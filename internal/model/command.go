package model

type CommandName string

const (
	CommandUnknown CommandName = "unknown"
	CommandStart   CommandName = "start"
	CommandHelp    CommandName = "help"
)

type Command struct {
	Name            CommandName
	ChatID          int64
	UserID          int64
	UserDisplayName string
	RawArgs         string
}
