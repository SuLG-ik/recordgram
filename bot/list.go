package bot

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"recordgram/botapi"
	"recordgram/records"
)

type ListCommand struct {
	db *gorm.DB
}

func (command ListCommand) Handler(ctx *botapi.MessageContext) {
	onlineRecords, err := records.FindNonMarked(command.db)
	if err != nil {
		log.WithError(err).Error("ListCommandHandler: error fetching onlineRecords")
	}
	if len(onlineRecords) == 0 {
		_, err := ctx.SendMessage("🗃 Список заявок пуст.")
		if err != nil {
			log.WithError(err).Error("ListCommandHandler: error sending message")
		}
	} else {
		_, err := ctx.SendMessage("🗃 Список заявок:")
		if err != nil {
			log.WithError(err).Error("ListCommandHandler: error sending message")
		}
	}
	for _, record := range onlineRecords {
		SendRecord(ctx.Bot, command.db, ctx.Update.Message.Chat.ID, record)
	}
}
