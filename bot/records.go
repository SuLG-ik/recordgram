package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"recordgram/botapi"
	"recordgram/messages"
	"recordgram/records"
	"strconv"
)

func SendRecord(bot *tgbotapi.BotAPI, db *gorm.DB, chatId int64, record records.OnlineRecord) {
	response, err := bot.Send(prepareRecordStatusMessage(chatId, record))
	if err != nil {
		log.WithError(err).Warnf("RecordsSender: error sending message")
	}
	_, err = messages.Create(db, chatId, int64(response.MessageID), record.Id)
	if err != nil {
		log.WithError(err).Errorf("CompanyCommandHandler: error saving to db")
	}
}

func createMarkupReplyButtons(record records.OnlineRecord) tgbotapi.InlineKeyboardMarkup {
	buttons := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔴 Отклонить", fmt.Sprintf("markdeny %v", record.Id)),
			tgbotapi.NewInlineKeyboardButtonData("🟢 Принять", fmt.Sprintf("markaccept %v", record.Id)),
		),
	)
	return buttons
}

type MarkAcceptHandler struct {
	db *gorm.DB
}

func (command MarkAcceptHandler) Handler(ctx *botapi.CallbackContext) {
	payload := ctx.Payload()
	number, err := strconv.ParseInt(payload, 10, 0)
	if err != nil {
		_, err := ctx.SendMessage("Ошибка. Сообщение устарело. Введите /list чтобы обновить список заявок")
		if err != nil {
			log.WithError(err).Warnf("MarkAcceptHandler: error sending message")
		}
		return
	}
	markRecord(command.db, ctx, number, true)
}

type MarkDenyHandler struct {
	db *gorm.DB
}

func (command MarkDenyHandler) Handler(ctx *botapi.CallbackContext) {
	payload := ctx.Payload()
	number, err := strconv.ParseInt(payload, 10, 0)
	if err != nil {
		_, err := ctx.SendMessage("Ошибка. Сообщение устарело. Введите /list чтобы обновить список заявок")
		if err != nil {
			log.WithError(err).Warnf("MarkAcceptHandler: error sending message")
		}
		return
	}
	markRecord(command.db, ctx, number, false)
}

func markRecord(db *gorm.DB, ctx *botapi.CallbackContext, recordId int64, marked bool) {
	record, err := records.FindById(db, recordId)
	if err != nil {
		_, err := ctx.SendMessage("Запись не найдена")
		if err != nil {
			log.WithError(err).Warnf("MarkAcceptHandler: error sending message")
		}
		return
	}
	err = records.Mark(db, &record, marked)
	if err != nil {
		log.WithError(err).Warnf("MarkAcceptHandler: error sending message")
		return
	}
	sendedMessages, err := messages.FindByRecordId(db, record.Id)
	if err != nil {
		log.WithError(err).Warnf("MarkAcceptHandler: error finding records")
		_, err := ctx.SendMessage("Ошибка обновления сообщений.")
		if err != nil {
			log.WithError(err).Warnf("MarkAcceptHandler: error sending message")
			return
		}
	}

	message := prepareRecordStatus(record)
	for _, sendedMessage := range sendedMessages {
		go EditMessages(sendedMessage, message, ctx)
	}
}

func EditMessages(message messages.RecordToMessage, text string, ctx *botapi.CallbackContext) {
	_, err := ctx.Bot.Send(tgbotapi.NewEditMessageText(message.ChatId, int(message.MessageId), text))
	if err != nil {
		log.WithError(err).Warnf("MarkAcceptHandler: error edit messages")
		_, err = ctx.SendMessage("Ошибка обновления сообщений.")
		if err != nil {
			log.WithError(err).Warnf("MarkAcceptHandler: error sending message")
			return
		}
	}
}

func prepareRecordStatus(record records.OnlineRecord) string {
	var status string
	if record.Marked == nil {
		status = "🟡 Статус: не отмечено"
	} else if *record.Marked {
		status = "🟢 Статус: принято"
	} else {
		status = "🔴 Статус: отклонено"
	}
	return fmt.Sprintf("📝Запись.\n%v\n💇 Имя: %v\n📲 Номер тeлeфона: +%v", status, record.Name, record.Phone)
}

func prepareRecordStatusMessage(chatId int64, record records.OnlineRecord) tgbotapi.MessageConfig {
	buttons := createMarkupReplyButtons(record)
	return NewMessageWithInlineKeyboardMarkup(chatId, prepareRecordStatus(record), buttons)
}
