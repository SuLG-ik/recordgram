package companies

import (
	"gorm.io/gorm"
	"recordgram/passwords"
)

type Company struct {
	Id        int64
	Title     string
	TokenHash string
	ChatId    int64
}

func FindById(db *gorm.DB, id int64) (Company, error) {
	company := Company{}
	err := db.Take(&company, id).Error
	return company, err
}
func FindByChatId(db *gorm.DB, id int64) (Company, error) {
	company := Company{}
	err := db.Take(&company, "chat_id = ?", id).Error
	return company, err
}

func CreateCompany(db *gorm.DB, chatId int64, title, token string) (Company, error) {
	hash, err := passwords.HashPassword(token)
	if err != nil {
		return Company{}, err
	}
	company := Company{ChatId: chatId, TokenHash: hash, Title: title}
	err = db.Create(&company).Error
	return company, err
}

func EditCompanyToken(db *gorm.DB, company *Company, token string) error {
	password, err := passwords.HashPassword(token)
	if err != nil {
		return err
	}
	company.TokenHash = password
	err = db.Save(&company).Error
	return err
}

func Delete(db *gorm.DB, company *Company) error {
	return db.Delete(&company).Error
}
