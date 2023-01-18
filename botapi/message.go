package botapi

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"recordgram/config"
	"strings"
)

type MessageHandlerFunc func(ctx *MessageContext)

type MessageContext struct {
	Update    tgbotapi.Update
	Bot       *tgbotapi.BotAPI
	Command   string
	ChatId    int64
	MessageId int
	Text      string
	Config    config.Config

	payload *string
}

func NewMessageContext(api *tgbotapi.BotAPI, update tgbotapi.Update, command string, config config.Config) MessageContext {
	return MessageContext{
		Update:    update,
		Bot:       api,
		Command:   command,
		ChatId:    update.Message.Chat.ID,
		MessageId: update.Message.MessageID,
		Text:      update.Message.Text,
		Config:    config,
	}
}

func (m *MessageContext) SendMessage(message string) (tgbotapi.Message, error) {
	return m.Bot.Send(tgbotapi.NewMessage(m.ChatId, message))
}
func (m *MessageContext) SendMessagef(message string, v ...any) (tgbotapi.Message, error) {
	return m.Bot.Send(tgbotapi.NewMessage(m.ChatId, fmt.Sprintf(message, v...)))
}

func (m *MessageContext) SendMarkdownMessage(message string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(m.ChatId, message)
	msg.ParseMode = "markdown"
	return m.Bot.Send(msg)
}
func (m *MessageContext) SendMarkdownMessagef(message string, v ...any) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(m.ChatId, fmt.Sprintf(message, v...))
	msg.ParseMode = "markdown"
	return m.Bot.Send(msg)
}

func (m *MessageContext) Payload() string {
	if m.payload != nil {
		return *m.payload
	}
	message := m.Update.Message.Text
	payload := strings.TrimLeft(strings.TrimPrefix(message, "/"+m.Command), " ")
	m.payload = &payload
	return payload
}
