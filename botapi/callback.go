package botapi

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type CallbackHandlerFunc func(ctx *CallbackContext)

type CallbackContext struct {
	Update tgbotapi.Update
	Bot    *tgbotapi.BotAPI

	payload *string
	command *string
}

func NewCallbackContext(api *tgbotapi.BotAPI, update tgbotapi.Update) CallbackContext {
	return CallbackContext{
		Update: update,
		Bot:    api,
	}
}

func (c *CallbackContext) Command() string {
	if c.command != nil {
		return *c.command
	}
	data := c.Update.CallbackData()
	i := strings.Index(data, " ")
	var command string
	if i > 0 {
		command = data[:i]
	} else {
		command = data
	}
	if i = strings.Index(command, "@"); i != -1 {
		command = command[:i]
	}
	c.command = &command
	return command
}

func (c *CallbackContext) Payload() string {
	if c.payload != nil {
		return *c.payload
	}
	message := c.Update.CallbackData()
	payload := strings.TrimLeft(strings.TrimPrefix(message, c.Command()), " ")
	c.payload = &payload
	return payload
}

func (c *CallbackContext) SendMessage(message string) (tgbotapi.Message, error) {
	return c.Bot.Send(tgbotapi.NewMessage(c.Update.CallbackQuery.From.ID, message))
}
func (c *CallbackContext) SendMessagef(message string, v ...any) (tgbotapi.Message, error) {
	return c.Bot.Send(tgbotapi.NewMessage(c.Update.CallbackQuery.From.ID, fmt.Sprintf(message, v...)))
}
