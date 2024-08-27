package db

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ttsBot/db/models"
	"ttsBot/logger"
	"ttsBot/types"
)

func AddMediaFile(client *mongo.Client, media *types.Media) {
	var file models.File

	file.Text = media.GetMessage()
	file.Filepath = media.GetPath()

	collection := client.Database("ttsBot").Collection("files")
	_, err := collection.InsertOne(context.TODO(), file)
	if err != nil {
		logger.Log.Warn(err)
		return
	}

	log.Info().Msg(fmt.Sprintf("Successfully added file to the database. %s ", file.Filepath))
}

func GetMediaFile(client *mongo.Client, text string) *types.Media {
	media := &types.Media{}
	var file models.File

	collection := client.Database("ttsBot").Collection("files")
	fmt.Println(text)
	filter := bson.D{{"text", text}}

	err := collection.FindOne(context.TODO(), filter).Decode(&file)

	if err != nil {
		logger.Log.Warn("!MONGODB! ", err)
		return nil
	}

	media.SetMessage(file.Text)
	media.SetPath(file.Filepath)

	return media
}
