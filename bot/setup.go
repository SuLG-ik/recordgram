package bot

import (
	"gorm.io/gorm"
	"recordgram/botapi"
	"recordgram/config"
)

func SetupBot(config config.Config, db *gorm.DB) *botapi.Bot {
	bot := botapi.NewBot(config)
	bot.HandleCommand("help", HelpCommand{}.Handler)
	bot.HandleCallback("markaccept", MarkAcceptHandler{db: db}.Handler)
	bot.HandleCallback("markdeny", MarkDenyHandler{db: db}.Handler)
	bot.HandleCommands([]string{"list", "l"}, ListCommand{db: db}.Handler)
	return bot
}
