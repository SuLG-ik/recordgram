package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func Init() {
	output, err := parseFlags()
	if err == nil {
		file, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil && file != nil {
			log.SetFormatter(&log.JSONFormatter{})
			log.SetOutput(file)
			log.WithField("path", output).Info("Loggers: initialized to file")
			return
		}
	}
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.Info("Loggers: initialized to console")
}

func parseFlags() (string, error) {
	logPath := os.Getenv("LOG")
	if err := validateConfigPath(logPath); err != nil {
		return "", err
	}
	return logPath, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
