package bot

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"recordgram/botapi"
	"recordgram/companies"
	"recordgram/utils"
)

type NewCompanyHandler struct {
	db *gorm.DB
}

func (command NewCompanyHandler) Handler(ctx *botapi.MessageContext) {
	_, err := companies.FindByChatId(command.db, ctx.ChatId)
	if err == nil {
		_, err := ctx.SendMessage("Компания уже привязана\n/disable")
		if err != nil {
			log.WithError(err).Warnf("NewCompanyCommandHandler: error sending message")
		}
		return
	}
	if len(ctx.Payload()) == 0 {
		_, err := ctx.SendMessage("Введите /new <название компании>")
		if err != nil {
			log.WithError(err).Warnf("NewCompanyCommandHandler: error sending message")
		}
		return
	}
	token := utils.GenerateTokenFromConfig(ctx.Config)
	company, err := companies.CreateCompany(command.db, ctx.ChatId, ctx.Payload(), token)
	if err != nil {
		_, err := ctx.SendMessage("Ошибка создания компании")
		if err != nil {
			log.WithError(err).Warnf("NewCompanyCommandHandler: error sending message")
		}
		return
	}
	_, err = ctx.SendMarkdownMessagef("Компания: *%v*\nAPI-токен: `%v:%v`", company.Title, company.Id, token)
	if err != nil {
		log.WithError(err).Warnf("NewCompanyCommandHandler: error sending message")
	}
}
