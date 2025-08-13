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

var a *app = &app{
	errLogger:  log.New(os.Stderr, "[ERROR]", log.Flags()),
	warnLogger: log.New(os.Stdout, "[WARN]", log.Flags()),
	infoLogger: log.New(os.Stdout, "[INFO]", log.Flags()),
}

func InitApp() *app {
	var err error
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	var logWriter io.Writer
	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		a.Warn("LOG_FILE not passed. Falling back to STDOUT")
		logWriter = os.Stdout
	} else {
		logWriter, err = os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			a.Fatal(fmt.Sprintf("Error opening file %s: %v\n", logFile, err))
		}
	}
	a.errLogger.SetOutput(logWriter)
	a.infoLogger.SetOutput(logWriter)
	a.warnLogger.SetOutput(logWriter)
	a.Port = port
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
