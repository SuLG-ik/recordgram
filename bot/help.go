package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"recordgram/botapi"
)

const (
	helpCommand = `📝Список доступных команд:
/help - инструкция
/list (/l) - список неомеченных заявок
/new - создание компании
/company - информация по компании
/disable - удаление клмпании
/key - перевыпуск Api-токена`
)

type HelpCommand struct {
}

func (command HelpCommand) Handler(ctx *botapi.MessageContext) {
	_, err := ctx.Bot.Send(tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, helpCommand))
	if err != nil {
		log.WithError(err).Warnf("HelpCommandHandler: error sending message")
	}
}
