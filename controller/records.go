package controller

import (
	"github.com/go-chi/chi/v5"
	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"recordgram/bot"
	"recordgram/botapi"
	"recordgram/companies"
	"recordgram/passwords"
	"recordgram/records"
	"recordgram/utils"
	"regexp"
	"strconv"
	"strings"
)

var phoneRegex = regexp.MustCompile("^7[0-9]{10}")

func RecordsRouter(db *gorm.DB, bot *botapi.Bot) func(router chi.Router) {
	return func(router chi.Router) {
		router.Post("/", ExternalRecordRouter(db, bot))
	}
}

var digitCheck = regexp.MustCompile(`[^0-9]`)

func FilterNums(text string) string {
	return digitCheck.ReplaceAllString(text, "")
}

func ExternalRecordRouter(db *gorm.DB, botClient *botapi.Bot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		token := query.Get("token")
		if !utils.IsTokenValid(token, botClient.Config) {
			http.Error(w, "Токен не указан", http.StatusBadRequest)
			return
		}
		id, apiToken := splitToken(token)
		company, err := companies.FindById(db, id)
		if err != nil || !passwords.MatchPassword(apiToken, company.TokenHash) {
			http.Error(w, "Компания не доступна", http.StatusBadRequest)
			return
		}
		phone := FilterNums(query.Get("phone"))
		if !phoneRegex.MatchString(phone) {
			http.Error(w, "Телефон не передан", http.StatusBadRequest)
			return
		}
		name := query.Get("name")
		nameLen := len(name)
		if nameLen > 32 || nameLen < 3 {
			http.Error(w, "Имя не передано и/или длина меньше 3/больше 32", http.StatusBadRequest)
			return
		}
		message := query.Get("message")
		messageLen := len(message)
		if messageLen > 200 {
			http.Error(w, "Слишком длинное сообщение", http.StatusBadRequest)
			return
		}
		record, err := records.Create(db, id, name, phone, message)
		if err != nil {
			logging.WithError(err).Errorln("Внутренняя ошибка создания записи")
			http.Error(w, "Внутренняя ошибка. Скоро починим", http.StatusInternalServerError)
		}
		go bot.SendRecord(botClient.BotApi, db, company.ChatId, record)
		w.WriteHeader(http.StatusOK)
	}

}

func splitToken(token string) (int64, string) {
	x := strings.Split(token, ":")
	id, err := strconv.ParseInt(x[0], 10, 64)
	if err != nil {
		panic(err)
	}
	return id, x[1]
}
