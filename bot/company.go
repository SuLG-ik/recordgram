package bot

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	botapi "recordgram/botapi"
	"recordgram/companies"
)

type CompanyHandler struct {
	db *gorm.DB
}

func (command CompanyHandler) Handler(ctx *botapi.MessageContext) {
	if len(ctx.Payload()) != 0 {
		_, err := ctx.SendMessage("Введите просто /company")
		if err != nil {
			log.WithError(err).Warnf("CompanyCommandHandler: error sending message")
		}
		return
	}
	company, err := companies.FindByChatId(command.db, ctx.Update.Message.Chat.ID)
	if err != nil {
		log.WithError(err).Error("CompanyCommandHandler: error find company")
		_, err := ctx.SendMessage("Компания не привязана.\n/new - создать компанию")
		if err != nil {
			log.WithError(err).Warnf("CompanyCommandHandler: error sending message")
		}
		return
	}
	_, err = ctx.SendMarkdownMessagef("Компания: *%v*\nКлюч: /key - перевыпустить ключ", company.Title)
	if err != nil {
		log.WithError(err).Warnf("CompanyCommandHandler: error sending message")
	}
}
