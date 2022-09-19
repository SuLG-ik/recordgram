package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func NewMessageWithInlineKeyboardMarkup(chatID int64, text string, markup tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig {
	config := tgbotapi.NewMessage(chatID, text)
	config.BaseChat.ReplyMarkup = markup
	return config
}

func NewMessageWithReplyKeyboardMarkup(chatID int64, text string, markup tgbotapi.ReplyKeyboardMarkup) tgbotapi.MessageConfig {
	config := tgbotapi.NewMessage(chatID, text)
	config.BaseChat.ReplyMarkup = markup
	return config
}
