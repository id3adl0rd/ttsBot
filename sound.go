package main

import (
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"os"
	"ttsBot/db"
	"ttsBot/logger"
	"ttsBot/types"
)

func Play(b *Bot, g *types.Guild, filePath, gID, channelID string) {
	vc, err := b.session.ChannelVoiceJoin(gID, channelID, false, true)
	if err != nil {
		logger.Log.Warn(err)
		return
	}

	owd, err := os.Getwd()

	if err != nil {
		logger.Log.Info(err)
	}

	for media := range g.GetMedia() {
		if !vc.Ready {
			vc.Disconnect()
			logger.Log.Info("Reconnecting...")

			vc, err = b.session.ChannelVoiceJoin(gID, channelID, false, true)
			if err != nil {
				logger.Log.Warn(err)
				return
			}
		}

		dgvoice.PlayAudioFile(vc, fmt.Sprintf("%s/%s", owd, media.GetPath()), make(chan bool))
	}

	vc.Disconnect()
	return
}

func CreateTTS(b *Bot, m *discordgo.MessageCreate) *types.Media {
	message := m.Content //strings.ToLower(m.Content)
	result := b.lru.Get(message)

	if result != nil {
		logger.Log.Info("Message was founded in cache")
		return result.(*types.Media)
	}

	var media *types.Media
	media = db.GetMediaFile(Client, message)
	if media != nil {
		logger.Log.Info("Message was founded in db. Restored into the cache")
		b.lru.Set(media.GetMessage(), media)

		return media
	}

	speech := htgotts.Speech{Folder: "audio", Language: voices.Russian, Handler: &handlers.Native{}}

	logger.Log.Info("Trying to create speech file")

	fileUUID := uuid.New().String()
	filePath, err := speech.CreateSpeechFile(message, fileUUID)

	if err != nil {
		logger.Log.Warn("Speech creating failed: ", err)
		return nil
	}

	logger.Log.Info("Speech created:", filePath)

	media = types.NewMedia(m.Content, filePath)
	b.lru.Set(message, media)
	db.AddMediaFile(Client, media)

	return media
}
