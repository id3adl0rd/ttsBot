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

func (s *Sound) Load(path string, session *discordgo.Session, gID, channelID string) {
	vc, err := session.ChannelVoiceJoin(gID, channelID, false, true)
	if err != nil {
		Log.Warn(err)
		return
	}

	owd, err := os.Getwd()

	if err != nil {
		Log.Info(err)
	}

	dgvoice.PlayAudioFile(vc, fmt.Sprintf("%s/%s", owd, path), make(chan bool))

	vc.Disconnect()
}
