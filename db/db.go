package db

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"ttsBot/config"
)

func Connect(dbcfg *config.DBConfig) *mongo.Client {
	clientoptions := options.Client().ApplyURI(fmt.Sprintf(dbcfg.Url, dbcfg.Port))
	client, err := mongo.Connect(context.TODO(), clientoptions)

	if err != nil {
		log.Fatal().Err(err).Msg("Connect MongoDB Failed")
		return nil
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal().Err(err).Msg("Connect MongoDB Failed")
		return nil
	}

	log.Info().Msg("Connect MongoDB Success")
	return client
}
