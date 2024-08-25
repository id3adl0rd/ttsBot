package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"ttsBot/config"
	"ttsBot/logger"
)

var (
	Log        *logger.Zerolog
	Enviroment *config.Config
)

func init() {
	Log = logger.NewZerolog()
	path, err := os.Getwd()

	if err != nil {
		Log.Error(err)
	}

	Enviroment, err = config.InitConfig(path)

	if err != nil {
		Log.Error(err)
	}
}

func main() {
	Log.Info("Starting bot...")

	bot := NewBot()
	bot.session.AddHandler(bot.MessageHandler)

	err := bot.session.Open()
	if err != nil {
		Log.Errorf("Error opening session: %v", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	Log.Info("Stopping bot...")
	bot.session.Close()
}

func test() {
	fmt.Println("testing")
}
