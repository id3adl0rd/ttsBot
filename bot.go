package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"net/url"
	"strings"
	"sync"
	"ttsBot/cache"
	"ttsBot/logger"
	"ttsBot/types"
)

type Bot struct {
	session    *discordgo.Session
	lru        *cache.LRU
	guilds     map[string]*types.Guild
	guildNames map[string]string
	mutex      sync.Mutex
}

func NewBot() *Bot {
	if Enviroment.Bot.Token == "" {
		logger.Log.Error("Error creating Discord session. Missing token.")
		return nil
	}

	bot, err := discordgo.New("Bot " + Enviroment.Bot.Token)

	if err != nil {
		logger.Log.Error("Error creating Discord session,", err)
		return nil
	}

	return &Bot{
		session:    bot,
		guilds:     make(map[string]*types.Guild),
		guildNames: make(map[string]string),
		lru:        cache.NewLru(Enviroment.Misc.CacheSize),
		mutex:      sync.Mutex{},
	}
}

func (b *Bot) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		logger.Log.Warn("Author compared with User")
		return
	}

	channel, err := b.getVoiceChannelByMessageID(m)
	if err != nil {
		logger.Log.Info(err)
		return
	}

	u, _ := url.ParseRequestURI(m.Content)

	if u != nil {
		logger.Log.Info("Url shoudn't be voiced")
		return
	}

	if len(m.Attachments) > 0 {
		logger.Log.Info("Attachments shoudn't be voiced")
		return
	}

	if strings.HasPrefix(m.Content, "!stop") {
		stop(b, channel)
		return
	}

	filePath := ""

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		filePath, err = addToQueue(b, channel, m)
		if err != nil {
			logger.Log.Info(err)
		}
	}()
	wg.Wait()

	if filePath == "" {
		logger.Log.Info("filePath is empty")
		return
	}

	if !b.GetGuild(channel).IsPlaying() {
		go func() {
			ttsRoutine(b, channel, filePath, m)
		}()
	}
}

func (b *Bot) AddGuild(gID string, g *types.Guild) {
	b.mutex.Lock()
	b.guilds[gID] = g
	b.mutex.Unlock()
}

func (b *Bot) RemoveGuild(guildID string) {
	b.mutex.Lock()
	b.guilds[guildID] = nil
	b.mutex.Unlock()
}

func (b *Bot) GetGuild(guildID string) *types.Guild {
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

func addToQueue(b *Bot, channel string, m *discordgo.MessageCreate) (string, error) {
	g := b.GetGuild(channel)
	if g != nil {
		if g.GetMedia() == nil {
			g.PrepareMediaChannel(int(Enviroment.Misc.QueueSize))
		} else {
			if g.IsQueueFull() == true {
				return "", errors.New("Queue is full")
			}
		}
	} else {
		g = types.NewGuild(channel)
		g.PrepareMediaChannel(int(Enviroment.Misc.QueueSize))
		b.AddGuild(channel, g)
	}

	media := CreateTTS(b, m)

	logger.Log.Infof("ttsMessage for %s is succesfuly created", channel)

	g.Enqueue(media)

	return media.GetPath(), nil
}

func ttsRoutine(b *Bot, channel, filePath string, m *discordgo.MessageCreate) error {
	g := b.GetGuild(channel)
	if g != nil {
		if g.IsQueueFull() == true {
			return errors.New("queue is full")
		}
	}
	g.SetIsPlaying(true)

	Play(b, g, filePath, m.GuildID, channel)

	return nil
}
