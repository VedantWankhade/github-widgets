package config

import (
	"fmt"
	"io"
	"log"
	"os"
)

type app struct {
	Port       string
	errLogger  *log.Logger
	warnLogger *log.Logger
	infoLogger *log.Logger
}

var a *app

func GetApp() *app {
	if a == nil {
		var err error
		port := os.Getenv("PORT")
		if port == "" {
			port = "80"
		}
		var logWriter io.Writer
		logFile := os.Getenv("LOG_FILE")
		if logFile == "" {
			logWriter = os.Stdout
		} else {
			logWriter, err = os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
		a = &app{
			Port:       port,
			infoLogger: log.New(logWriter, "[INFO]", log.Flags()),
			warnLogger: log.New(logWriter, "[WARN]", log.Flags()),
			errLogger:  log.New(logWriter, "[ERROR]", log.Flags()),
		}
	}
	return a
}

func (a *app) GetErrLogger() *log.Logger {
	return a.errLogger
}

func (a *app) Info(msg string) {
	a.infoLogger.Println(msg)
}

func (a *app) Warn(msg string) {
	a.warnLogger.Println(msg)
}

func (a *app) Error(msg string) {
	a.errLogger.Println(msg)
}

func (a *app) Fatal(data ...any) {
	a.errLogger.Fatal(fmt.Sprintf("%v\nQuitting...\n", data))
}
