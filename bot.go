package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"sync"
	"ttsBot/cache"
)

type Bot struct {
	session    *discordgo.Session
	lru        *cache.LRU
	guilds     map[string]*Guild
	guildNames map[string]string
	mutex      sync.Mutex
}

func NewBot() *Bot {
	if Enviroment.Bot.Token == "" {
		Log.Error("Error creating Discord session. Missing token.")
		return nil
	}

	bot, err := discordgo.New("Bot " + Enviroment.Bot.Token)

	if err != nil {
		Log.Error("Error creating Discord session,", err)
		return nil
	}

	return &Bot{
		session:    bot,
		guilds:     make(map[string]*Guild),
		guildNames: make(map[string]string),
		lru:        cache.NewLru(100),
		mutex:      sync.Mutex{},
	}
}

func (b *Bot) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		Log.Warn("Author compared with User")
		return
	}

	channel, err := b.getVoiceChannelByMessageID(m)
	if err != nil {
		Log.Info(err)
		return
	}

	fmt.Println("lol", strings.HasPrefix(m.Content, "!stop"))
	if strings.HasPrefix(m.Content, "!stop") {
		stop(b, channel)
		return
	}

	go func() {
		ttsRoutine(b, channel, m)
	}()
}

func (b *Bot) AddGuild(m *discordgo.MessageCreate) {
	b.mutex.Lock()
	guild := NewGuild(b.getGuildName(m))
	b.guilds[m.GuildID] = guild
	b.mutex.Unlock()
}

func (b *Bot) AddGuildDirectly(gID string, g *Guild) {
	b.mutex.Lock()
	b.guilds[gID] = g
	b.mutex.Unlock()
}

func (b *Bot) RemoveGuild(guildID string) {
	b.mutex.Lock()
	b.guilds[guildID] = nil
	b.mutex.Unlock()
}

func (b *Bot) GetGuild(guildID string) *Guild {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.guilds[guildID]
}

func (b *Bot) getGuildName(message *discordgo.MessageCreate) string {
	value, found := b.guildNames[message.GuildID]
	if !found {
		guild, err := b.session.Guild(message.GuildID)
		if err != nil {
			b.guildNames[message.GuildID] = message.GuildID
			return message.GuildID
		}
		b.guildNames[message.GuildID] = guild.Name
		return guild.Name
	}

	return value
}

func (b *Bot) getVoiceChannelByMessageID(message *discordgo.MessageCreate) (string, error) {
	guild, err := b.session.State.Guild(message.GuildID)
	if err != nil {
		return "", err
	}

	for _, voiceStates := range guild.VoiceStates {
		if voiceStates.UserID == message.Author.ID {
			return voiceStates.ChannelID, nil
		}
	}

	return "", errors.New("user not in voice channel")
}

func (b *Bot) getVoiceChannelByMessageID2(message *discordgo.MessageCreate) (string, error) {
	guild, err := b.session.State.Guild(message.GuildID)
	if err != nil {
		return "", err
	}

	for _, voiceStates := range guild.VoiceStates {
		if voiceStates.UserID == message.Author.ID {
			return voiceStates.ChannelID, nil
		}
	}

	return "", errors.New("user not in voice channel")
}

func stop(b *Bot, channel string) {
	guild := b.GetGuild(channel)
	if guild != nil {
		guild.Stop()
	} else {
		b.session.Lock()
		if b.session.VoiceConnections[channel] != nil {
			b.session.VoiceConnections[channel].Disconnect()
		}
		b.session.Unlock()
	}
}

func ttsRoutine(b *Bot, channel string, m *discordgo.MessageCreate) error {
	g := b.GetGuild(channel)
	if g != nil {
		if g.IsQueueFull() == true {
			return errors.New("Queue is full")
		}
	} else {
		g = NewGuild(channel)
		b.AddGuild(m)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	filePath := CreateTTS(b, m)
	wg.Done()

	wg.Wait()

	Log.Infof("ttsMessage for %s is succesfuly created", channel)

	g.PrepareMediaChannel(10)

	g.Enqueue(NewMedia(m.Content, filePath))

	wg.Add(1)
	Play(b, g, filePath, m.GuildID, channel)
	wg.Done()

	wg.Done()

	return nil
}
