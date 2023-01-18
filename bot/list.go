package bot

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"recordgram/botapi"
	"recordgram/companies"
	"recordgram/records"
)

type ListCommand struct {
	db *gorm.DB
}

func (command ListCommand) Handler(ctx *botapi.MessageContext) {
	company, err := companies.FindByChatId(command.db, ctx.ChatId)
	if err != nil {
		_, err := ctx.SendMessage("Компания не привязана")
		if err != nil {
			log.WithError(err).Error("ListCommandHandler: error sending message")
		}
		return
	}
	onlineRecords, err := records.FindNonMarked(command.db, company.Id)
	if err != nil {
		log.WithError(err).Error("ListCommandHandler: error fetching onlineRecords")
	}
	if len(onlineRecords) == 0 {
		_, err := ctx.SendMessage("🗃 Список заявок пуст.")
		if err != nil {
			log.WithError(err).Error("ListCommandHandler: error sending message")
		}
		return
	}
	_, err = ctx.SendMessage("🗃 Список заявок:")
	if err != nil {
		log.WithError(err).Error("ListCommandHandler: error sending message")
	}
	for _, record := range onlineRecords {
		SendRecord(ctx.Bot, command.db, ctx.Update.Message.Chat.ID, record)
	}
}
