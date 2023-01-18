package records

import (
	"gorm.io/gorm"
	"time"
)

type OnlineRecord struct {
	Id        int64
	CompanyId int64
	Name      string
	Phone     string
	Date      time.Time
	Message   string
	Marked    *bool
}

func Create(db *gorm.DB, companyId int64, name string, phone string, message string) (OnlineRecord, error) {
	record := OnlineRecord{Name: name, Phone: phone, Date: time.Now(), Message: message, CompanyId: companyId}
	err := db.Create(&record).Error
	return record, err
}

func FindById(db *gorm.DB, markId int64) (OnlineRecord, error) {
	record := OnlineRecord{}
	err := db.Take(&record, markId).Error
	return record, err
}

func FindNonMarked(db *gorm.DB, companyId int64) ([]OnlineRecord, error) {
	var records []OnlineRecord
	err := db.Where("marked is NULL").Find(&records, "company_id = ?", companyId).Error
	return records, err
}

func Mark(db *gorm.DB, record *OnlineRecord, marked bool) error {
	record.Marked = &marked
	err := db.Save(&record).Error
	return err
}
