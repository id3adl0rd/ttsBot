package main

import "github.com/bwmarrin/discordgo"

type Bot struct {
	session *discordgo.Session
}

func NewBot() *Bot {
	if Enviroment.Bot.Token == "" {
		Log.Error("Error creating Discord session. Missing token.")
		return nil
	}

	dg, err := discordgo.New("Bot " + Enviroment.Bot.Token)

	if err != nil {
		Log.Error("Error creating Discord session,", err)
		return nil
	}

	return &Bot{
		dg,
	}
}
