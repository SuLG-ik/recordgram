package bot

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"recordgram/botapi"
	"recordgram/companies"
	"recordgram/utils"
)

type KeyHandler struct {
	db *gorm.DB
}

func (command KeyHandler) Handler(ctx *botapi.MessageContext) {
	company, err := companies.FindByChatId(command.db, ctx.ChatId)
	if err != nil {
		_, err := ctx.SendMessage("Компания не привязана\n/new <название компании> - создание компании")
		if err != nil {
			log.WithError(err).Warnf("KeyHandler: error sending message")
		}
		return
	}
	if len(ctx.Payload()) != 0 {
		_, err := ctx.SendMessage("Введите просто /key")
		if err != nil {
			log.WithError(err).Warnf("KeyHandler: error sending message")
		}
		return
	}
	token := utils.GenerateTokenFromConfig(ctx.Config)
	err = companies.EditCompanyToken(command.db, &company, token)
	if err != nil {
		_, err := ctx.SendMessage("Ошибка создания токена")
		if err != nil {
			log.WithError(err).Warnf("KeyHandler: error sending message")
		}
		return
	}
	_, err = ctx.SendMarkdownMessagef("Выпущен новый Api-токен\nКомпания: *%v*\nAPI-токен: `%v:%v`", company.Title, company.Id, token)
	if err != nil {
		log.WithError(err).Warnf("KeyHandler: error sending message")
	}
}
