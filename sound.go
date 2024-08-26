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
)

func Play(b *Bot, g *Guild, filePath, gID, channelID string) {
	vc, err := b.session.ChannelVoiceJoin(gID, channelID, false, true)
	if err != nil {
		Log.Warn(err)
		return
	}

	owd, err := os.Getwd()

	if err != nil {
		Log.Info(err)
	}

	for media := range g.media {
		if !vc.Ready {
			vc.Disconnect()
			Log.Info("Reconnecting...")

			vc, err = b.session.ChannelVoiceJoin(gID, channelID, false, true)
			if err != nil {
				Log.Warn(err)
				return
			}
		}

		dgvoice.PlayAudioFile(vc, fmt.Sprintf("%s/%s", owd, media.path), make(chan bool))
	}

	fmt.Println("lolzxc")
	g.UpdateTime()

	vc.Disconnect()
	return
}

func CreateTTS(b *Bot, m *discordgo.MessageCreate) string {
	message := m.Content
	result := b.lru.Get(message)

	if result != nil {
		Log.Info("Message was founded in cache")
		return result.(string)
	}

	speech := htgotts.Speech{Folder: "audio", Language: voices.Russian, Handler: &handlers.Native{}}

	Log.Info("Trying to create speech file")

	fileUUID := uuid.New().String()
	filePath, err := speech.CreateSpeechFile(message, fileUUID)

	if err != nil {
		Log.Warn("Speech creating failed: ", err)
		return ""
	}

	Log.Info("Speech created:", filePath)

	b.lru.Set(message, filePath)

	return filePath
}
