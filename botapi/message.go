package botapi

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type MessageHandlerFunc func(ctx *MessageContext)

type MessageContext struct {
	Update  tgbotapi.Update
	Bot     *tgbotapi.BotAPI
	Command string

	payload *string
}

func NewMessageContext(api *tgbotapi.BotAPI, update tgbotapi.Update, command string) MessageContext {
	return MessageContext{
		Update:  update,
		Bot:     api,
		Command: command,
	}
}

func (m *MessageContext) SendMessage(message string) (tgbotapi.Message, error) {
	return m.Bot.Send(tgbotapi.NewMessage(m.Update.Message.Chat.ID, message))
}
func (m *MessageContext) SendMessagef(message string, v ...any) (tgbotapi.Message, error) {
	return m.Bot.Send(tgbotapi.NewMessage(m.Update.Message.Chat.ID, fmt.Sprintf(message, v...)))
}

func (m *MessageContext) Payload() string {
	if m.payload != nil {
		return *m.payload
	}

	index := strings.Index(m.Update.Message.Text, m.Command)
	if len(m.Update.Message.Text) == index+1 {
		return ""
	}
	payload := m.Update.Message.Text[index:]
	m.payload = &payload
	return payload
}
