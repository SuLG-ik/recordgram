package db

import (
	"fmt"
	logging "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"recordgram/config"
	"recordgram/messages"
	"recordgram/records"
)

func migrate(config config.Config, conn *gorm.DB) {
	if !config.Datasource.AutoMigration {
		logging.Info("Database: automigration skipped")
		return
	}
	logging.Info("Database: automigrating")
	err := conn.AutoMigrate(&messages.RecordToMessage{}, &records.OnlineRecord{})
	if err != nil {
		logging.Panic(err)
	}
	logging.Info("Database: automigrated")
}

func SetupDB(config config.Config) *gorm.DB {
	db := connect(config)
	migrate(config, db)
	return db
}

func connect(config config.Config) *gorm.DB {
	datasource := config.Datasource
	var sslmode string
	if datasource.Ssl {
		sslmode = "enable"
	} else {
		sslmode = "disable"
	}
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		datasource.Host,
		datasource.User,
		datasource.Password,
		datasource.Database,
		datasource.Port,
		sslmode,
		datasource.Timezone,
	)
	conn := postgres.New(postgres.Config{DSN: dsn})
	logging.Info("Database: connecting")
	db, err := gorm.Open(conn)

	if err != nil {
		logging.WithError(err).Panicf("Database: error connecting")
	}
	logging.Info("Database: connected")
	return db
}
