package main

import (
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

	telegramBot := bot.SetupBot(loadedConfig, database)

	server := controller.NewServer(loadedConfig, database, telegramBot)

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()
		telegramBot.Start()
	}()
	go func() {
		defer waitGroup.Done()
		server.Start()
	}()

	log.Infof("Startup took aprox %v", time.Since(start))
	waitGroup.Wait()
}
