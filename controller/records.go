package controller

import (
	"github.com/go-chi/chi/v5"
	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"recordgram/bot"
	"recordgram/botapi"
	"recordgram/config"
	"recordgram/records"
	"regexp"
)

var phoneRegex, _ = regexp.Compile("^7[0-9]{10}")

//var tokenRegex, _ = regexp.Compile("[a-zA-Z0-9]{32}")

func RecordsRouter(config config.Config, db *gorm.DB, bot *botapi.Bot) func(router chi.Router) {
	return func(router chi.Router) {
		router.Post("/", ExternalRecordRouter(config, db, bot))
	}
}

func ExternalRecordRouter(config config.Config, db *gorm.DB, botClient *botapi.Bot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		phone := query.Get("phone")
		if !phoneRegex.MatchString(phone) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		name := query.Get("name")
		nameLen := len(name)
		if nameLen > 32 || nameLen < 3 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		record, err := records.Create(db, name, phone)
		if err != nil {
			logging.WithError(err).Errorln("Внутренняя ошибка создания записи")
			http.Error(w, "Внутренняя ошибка. Скоро починим", http.StatusInternalServerError)
		}
		go bot.SendRecord(botClient.BotApi, db, config.Telegram.ChatId, record)
		w.WriteHeader(http.StatusOK)
	}

}
