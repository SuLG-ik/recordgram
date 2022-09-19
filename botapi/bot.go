package botapi

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"recordgram/config"
	"strings"
	"sync/atomic"
)

//type MessageHandlerArgs struct {
//	update  tgbotapi.Update
//	command string
//	args    []string
//}
//
//type MessageHandler struct {
//	name    string
//	alias   []string
//	handler func(botapi *tgbotapi.BotAPI, args MessageHandlerArgs, deps *deps.Deps)
//}

type Bot struct {
	BotApi           *tgbotapi.BotAPI
	messageHandlers  map[string]MessageHandlerFunc
	callbackHandlers map[string]CallbackHandlerFunc
	config           config.Config
	started          atomic.Bool
}

func (bot *Bot) HandleCommand(endpoint string, handlerFunc MessageHandlerFunc) {
	if bot.started.Load() {
		log.Panicf("TelegramBot: already started. Adding command handler canceled")
	}
	command := strings.TrimPrefix(endpoint, "/")
	_, exists := bot.messageHandlers[command]
	if exists {
		log.Panicf("TelegramBot: command handler already exists [%v]", command)
	}
	bot.messageHandlers[command] = handlerFunc
}
func (bot *Bot) HandleCommands(endpoints []string, handlerFunc MessageHandlerFunc) {
	for _, endpoint := range endpoints {
		bot.HandleCommand(endpoint, handlerFunc)
	}
}

func (bot *Bot) HandleCallback(endpoint string, handlerFunc CallbackHandlerFunc) {
	if bot.started.Load() {
		log.Panicf("TelegramBot: already started. Adding command handler canceled")
	}
	command := strings.TrimPrefix(endpoint, "/")
	_, exists := bot.messageHandlers[command]
	if exists {
		log.Panicf("TelegramBot: command handler already exists [%v]", command)
	}
	bot.callbackHandlers[command] = handlerFunc
}
func (bot *Bot) HandleCallbacks(endpoints []string, handlerFunc CallbackHandlerFunc) {
	for _, endpoint := range endpoints {
		bot.HandleCallback(endpoint, handlerFunc)
	}
}

func NewBot(config config.Config) *Bot {
	log.Info("TelegramBot: initializing")
	botApi, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		log.WithError(err).Panic("TelegramBot: error bot initializing")
	}
	err = tgbotapi.SetLogger(log.StandardLogger())
	if err != nil {
		log.WithError(err).Warnf("TelegramBot: error attach logger")
	}
	log.WithField("username", botApi.Self.UserName).Infof("TelegramBot: initialized")
	bot := &Bot{
		BotApi:           botApi,
		messageHandlers:  map[string]MessageHandlerFunc{},
		callbackHandlers: map[string]CallbackHandlerFunc{},
		config:           config,
	}
	return bot
}

//func NewBot(deps *deps.Deps) *Bot {
//	err := tgbotapi.SetLogger(log.InfoLogger)
//	if err != nil {
//		log.Warn("Error attach telegram logger to info log", err)
//	}
//	config := deps.Config
//	log.Info("TelegramBot: starting")
//	BotApi, err := tgbotapi.NewBotAPI(config.Telegram.Token)
//	deps.Bot = BotApi
//	if err != nil {
//		log.Fatal(err)
//	}
//	me, err := BotApi.GetMe()
//	if err != nil {
//		log.Panic(errors.New("can't get me"))
//	}
//	log.Info("TelegramBot: started with username", me.UserName)
//	BotApi.Debug = config.Debug
//	handlers := toMap(Handlers)
//	if _, ok := handlers["help"]; ok {
//		panic("multiple handler initialization: help")
//	}
//	handlers["help"] = HelpCommand
//	botapi := Bot{BotApi: BotApi, preloadedHandlers: handlers, deps: deps}
//	return &botapi
//}

//func toMap(handlers []MessageHandler) map[string]MessageHandler {
//	result := map[string]MessageHandler{}
//	for _, handler := range handlers {
//		if _, ok := result[handler.name]; ok {
//			panic("multiple handler initialization: " + handler.name)
//		}
//		result[handler.name] = handler
//		for _, alias := range handler.alias {
//			result[alias] = handler
//		}
//	}
//	return result
//}

func (bot *Bot) Start() {
	if bot.started.Load() {
		log.Panicf("TelegramBot: already started. Starting canceled")
	}
	bot.started.Store(true)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage, tgbotapi.UpdateTypeCallbackQuery}

	updates := bot.BotApi.GetUpdatesChan(u)

	log.WithField("handlers", lo.Keys(bot.messageHandlers)).Debug("TelegramBot: started with handlers")
	for update := range updates {
		go produceUpdate(bot, update)
	}
}

func produceUpdate(bot *Bot, update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{"err": err, "update": update.UpdateID}).Warnf("TelegramBot: error handling update")
		}
	}()
	if produceMessage(bot, update) {
		return
	}
	if produceCallback(bot, update) {
		return
	}

	log.Warnf("TelegramBot: Update does not send to user %v", update.UpdateID)
}

func produceMessage(bot *Bot, update tgbotapi.Update) bool {
	if update.Message != nil {
		text := update.Message.Text
		if strings.HasPrefix(text, "/") {
			executeCommand(bot, update)
			return true
		}
		_, err := bot.BotApi.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "/help - список команд"))
		if err != nil {
			log.Warnf("TelegramBot: Update does not send to user %v", update.Message.Chat.ID)
		}
		return true
	}
	return false
}

func produceCallback(bot *Bot, update tgbotapi.Update) bool {
	if update.CallbackQuery != nil {
		executeCallback(bot, update)
		return true
	}
	return false
}

func executeCommand(bot *Bot, update tgbotapi.Update) {
	command := update.Message.Command()
	if len(command) == 0 {
		_, _ = bot.BotApi.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "/help - список команд"))
	}
	handler, ok := bot.messageHandlers[command]
	if ok {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{"command": command, "chat_id": update.Message.Chat.ID, "message": update.Message.Text, "error": err}).Warnf("TelegramBot: error handling command")
			}
		}()
		context := NewMessageContext(bot.BotApi, update, command)
		handler(&context)
		return
	}
	_, _ = bot.BotApi.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "/help - список команд"))
}

func executeCallback(bot *Bot, update tgbotapi.Update) {
	ctx := NewCallbackContext(bot.BotApi, update)
	handler, ok := bot.callbackHandlers[ctx.Command()]
	if ok {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{"command": ctx.Command(), "chat_id": update.CallbackQuery.From.ID, "message": update.Message.Text, "error": err}).Warnf("TelegramBot: error handling command")
			}
		}()
		handler(&ctx)
		return
	}
}
