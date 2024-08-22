package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"strings"
	"sync"
)

type Bot struct {
	Session *discordgo.Session
	Cache   *LRU
	Mutex   sync.Mutex
	Wg      sync.WaitGroup
	Queue   map[int]string
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
		Session: dg,
		Cache:   NewLru(100),
		Mutex:   sync.Mutex{},
		Wg:      sync.WaitGroup{},
		Queue:   make(map[int]string),
	}
}

func (b *Bot) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		Log.Warn("Author compared with User")
		return
	}

	if strings.HasPrefix(m.Content, "!test") {

	}

	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		Log.Warn(err)
		return
	}

	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		Log.Warn(err)
		return
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			Log.Info("Joining voice channel")

			b.PlayTTS(m.Content, g.ID, vs.ChannelID, m)

			Log.Info("Successfully played")

			return
		}
	}
}

func (b *Bot) CreateTTS(message string) string {
	result := b.Cache.Get(message)

	if result != nil {
		Log.Info("Message was founded in cache")
		return result.(string)
	}

	speech := htgotts.Speech{Folder: "audio", Language: voices.Russian, Handler: &handlers.Native{}}

	Log.Info("Trying to create speech file")

	fileUUID := uuid.New().String()
	fileName, err := speech.CreateSpeechFile(message, fileUUID)

	if err != nil {
		Log.Warn("Speech creating failed: ", err)
		return ""
	}

	Log.Info("Speech created:", fileName)

	b.Cache.Set(message, fileName)

	return fileName
}

func (b *Bot) PlayTTS(message, gID, vsChannelID string, m *discordgo.MessageCreate) {
	var filePath string

	b.Wg.Add(1)
	go func() {
		filePath = b.CreateTTS(message)
		b.Wg.Done()
	}()

	b.Wg.Wait()

	if filePath == "" {
		Log.Warn("TTS sound doesn't created")
		b.Session.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		return
	}

	b.Session.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	b.Wg.Add(1)
	go func() {
		sound := NewSound(filePath)
		sound.Load(filePath, b.Session, gID, vsChannelID)

		b.Wg.Done()
	}()

	b.Wg.Wait()
}

func (b *Bot) AddToQueue(message string) {

}

func (b *Bot) RemoveFromQueue(message string) {

}

func (b *Bot) Stop() {

}
