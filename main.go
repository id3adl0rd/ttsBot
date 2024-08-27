package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"ttsBot/config"
	"ttsBot/db"
	"ttsBot/logger"
)

var (
	Enviroment *config.Config
	Client     *mongo.Client
)

func init() {
	logger.Log = logger.NewZerolog()
	path, err := os.Getwd()

	if err != nil {
		logger.Log.Error(err)
	}

	Enviroment, err = config.InitConfig(path)

	if err != nil {
		logger.Log.Error(err)
	}

	Client = db.Connect(Enviroment.DBConfig)

	if Client == nil {
		return
	}
}

func main() {
	logger.Log.Info("Starting bot...")

	bot := NewBot()

	runtime.GOMAXPROCS(6)

	bot.session.AddHandler(bot.MessageHandler)

	err := bot.session.Open()
	if err != nil {
		logger.Log.Errorf("Error opening session: %v", err)
	}

	go clearConnection(bot)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	logger.Log.Info("Stopping bot...")
	bot.session.Close()

	logger.Log.Info("Stopping database...")
	Client.Disconnect(context.Background())
}
