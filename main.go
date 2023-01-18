package main

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"recordgram/bot"
	"recordgram/config"
	"recordgram/controller"
	"recordgram/db"
	logging2 "recordgram/logging"
	"sync"
	"time"
)

var start = time.Now()

func main() {
	logging2.Init()
	loadedConfig := config.FromFile()
	database := db.SetupDB(loadedConfig)
	sqlDb, _ := database.DB()
	defer func(sqlDb *sql.DB) {
		err := sqlDb.Close()
		if err != nil {
			log.WithError(err).Infof("DB: close")
		}
	}(sqlDb)

	telegramBot := bot.SetupBot(loadedConfig, database)

	server := controller.NewServer(loadedConfig, database, telegramBot)

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		telegramBot.Start(loadedConfig)
	}()
	go func() {
		defer waitGroup.Done()
		server.Start()
	}()

	log.Infof("Startup took aprox %v", time.Since(start))
	waitGroup.Wait()
}
