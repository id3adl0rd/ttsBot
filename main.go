package main

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	Log        *Zerolog
	Enviroment *Config
	DBot       *Bot
)

const (
	LogFile = "ttsBot.log"
)

func init() {
	Log = NewZerolog()
	path, err := os.Getwd()

	if err != nil {
		Log.Error(err)
	}

	Enviroment, err = InitConfig(path)

	if err != nil {
		Log.Error(err)
	}
}

func main() {
	Log.Info("Starting bot...")

	bot := NewBot()
	DBot = bot

	bot.Session.AddHandler(bot.MessageHandler)

	err := bot.Session.Open()
	if err != nil {
		Log.Errorf("Error opening session: %v", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	Log.Info("Stopping bot...")
	bot.Session.Close()
}
