package bot

import (
	"gorm.io/gorm"
	"recordgram/botapi"
	"recordgram/config"
)

func SetupBot(config config.Config, db *gorm.DB) *botapi.Bot {
	bot := botapi.NewBot(config)
	bot.HandleCommand("help", HelpCommand{}.Handler)
	bot.HandleCommand("company", CompanyHandler{db: db}.Handler)
	bot.HandleCommand("new", NewCompanyHandler{db: db}.Handler)
	bot.HandleCommand("disable", DisableCompanyHandler{db: db}.Handler)
	bot.HandleCommand("disableagree", DisableAgreeCompanyHandler{db: db}.Handler)
	bot.HandleCommand("key", KeyHandler{db: db}.Handler)
	bot.HandleCommands([]string{"list", "l"}, ListCommand{db: db}.Handler)
	bot.HandleCallback("markaccept", MarkAcceptHandler{db: db}.Handler)
	bot.HandleCallback("markdeny", MarkDenyHandler{db: db}.Handler)
	return bot
}
