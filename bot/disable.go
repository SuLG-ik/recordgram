package bot

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"recordgram/botapi"
	"recordgram/companies"
)

type DisableCompanyHandler struct {
	db *gorm.DB
}

func (command DisableCompanyHandler) Handler(ctx *botapi.MessageContext) {
	company, err := companies.FindByChatId(command.db, ctx.ChatId)
	if err != nil {
		_, err := ctx.SendMessage("Компания ещё не привязана\n/new - создание компании")
		if err != nil {
			log.WithError(err).Warnf("DisableCompanyHandler: error sending message")
		}
		return
	}

	_, err = ctx.SendMarkdownMessagef("Подтверждение удаления компании: *%v*\n/disableagree", company.Title)
	if err != nil {
		log.WithError(err).Warnf("DisableCompanyHandler: error sending message")
	}
}

type DisableAgreeCompanyHandler struct {
	db *gorm.DB
}

func (command DisableAgreeCompanyHandler) Handler(ctx *botapi.MessageContext) {
	company, err := companies.FindByChatId(command.db, ctx.ChatId)
	if err != nil {
		_, err := ctx.SendMessage("Компания ещё не привязана\n/new <название компании> - создание компании")
		if err != nil {
			log.WithError(err).Warnf("DisableCompanyHandler: error sending message")
		}
		return
	}
	err = companies.Delete(command.db, &company)
	if err != nil {
		_, err := ctx.SendMessage("Компания не удалена\n")
		log.WithError(err).Errorf("DisableCompanyHandler: error deleting company")
	}
	_, err = ctx.SendMarkdownMessagef("Компания удалена: *%v*\n/new <название компании> - создание компании", company.Title)
	if err != nil {
		log.WithError(err).Errorf("DisableCompanyHandler: error sending ьуыыфпу")
	}
}
