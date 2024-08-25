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

/*
type ISound interface {
	Play()
	Load()
}

type Sound struct {
	Path string
}

func NewSound(path string) *Sound {
	return &Sound{path}
}

func (s *Sound) Play(name string, vc *discordgo.VoiceConnection) {
	//s.Load(name, vc)
}

func (b *Bot) Load(path string, session *discordgo.Session, gID, channelID string) {
	if b.voiceChannel == nil {
		vc, err := session.ChannelVoiceJoin(gID, channelID, false, true)
		if err != nil {
			Log.Warn(err)
			return
		}

		b.voiceChannel = vc
	}

	owd, err := os.Getwd()

	if err != nil {
		Log.Info(err)
	}

	stop := make(chan bool)
	fmt.Println("Playing playing!!!")
	dgvoice.PlayAudioFile(b.voiceChannel, fmt.Sprintf("%s/%s", owd, path), stop)
}*/

func playSound(b *Bot, filePath, gID, channelID string) {
	vc, err := b.session.ChannelVoiceJoin(gID, channelID, false, true)
	if err != nil {
		Log.Warn(err)
		return
	}

	owd, err := os.Getwd()

	if err != nil {
		Log.Info(err)
	}

	dgvoice.PlayAudioFile(vc, fmt.Sprintf("%s/%s", owd, filePath), make(chan bool))
}

func createTTS(b *Bot, m *discordgo.MessageCreate) string {
	message := m.Content
	/*result := b.cache.Get(message)

	if result != nil {
		Log.Info("Message was founded in cache")
		return result.(string)
	}*/

	speech := htgotts.Speech{Folder: "audio", Language: voices.Russian, Handler: &handlers.Native{}}

	Log.Info("Trying to create speech file")

	fileUUID := uuid.New().String()
	filePath, err := speech.CreateSpeechFile(message, fileUUID)

	if err != nil {
		Log.Warn("Speech creating failed: ", err)
		return ""
	}

	Log.Info("Speech created:", filePath)

	//b.cache.Set(message, filePath)

	return filePath
}
