package records

import (
	"gorm.io/gorm"
)

type OnlineRecord struct {
	Id     int64
	Name   string
	Phone  string
	Marked *bool
}

func Create(db *gorm.DB, name string, phone string) (OnlineRecord, error) {
	record := OnlineRecord{Name: name, Phone: phone}
	err := db.Create(&record).Error
	return record, err
}

func FindById(db *gorm.DB, markId int64) (OnlineRecord, error) {
	record := OnlineRecord{}
	err := db.Take(&record, markId).Error
	return record, err
}

func FindNonMarked(db *gorm.DB) ([]OnlineRecord, error) {
	var records []OnlineRecord
	err := db.Where("marked is NULL").Find(&records).Error
	return records, err
}

func Mark(db *gorm.DB, record *OnlineRecord, marked bool) error {
	record.Marked = &marked
	err := db.Save(&record).Error
	return err
}
