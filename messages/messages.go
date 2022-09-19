package messages

import (
	"gorm.io/gorm"
)

type RecordToMessage struct {
	Id        int64
	ChatId    int64
	MessageId int64
	RecordId  int64
}

func FindByRecordId(db *gorm.DB, recordId int64) ([]RecordToMessage, error) {
	var record []RecordToMessage
	err := db.Where("record_id = ?", recordId).Find(&record).Error
	return record, err
}

func FindByMessageId(db *gorm.DB, chatId int64, messageId int64) (RecordToMessage, error) {
	record := RecordToMessage{}
	err := db.Take(&record, "chat_id = ?", chatId, "message_id", messageId).Error
	return record, err
}

func Create(db *gorm.DB, chatId int64, messageId int64, recordId int64) (RecordToMessage, error) {
	record := RecordToMessage{ChatId: chatId, MessageId: messageId, RecordId: recordId}
	err := db.Create(&record).Error
	return record, err
}
