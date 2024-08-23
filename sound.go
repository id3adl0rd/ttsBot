package main

import (
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"os"
)

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
}
