package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"strings"
	"sync"
	"time"
)

type Bot struct {
	session      *discordgo.Session
	voiceChannel *discordgo.VoiceConnection
	cache        *LRU
	mutex        sync.RWMutex
	wg           sync.WaitGroup
	queue        []string
	channelID    string
	guildID      string
	isPlaying    bool
}

type Guild struct {
	channelID string
	guildID   string
	isPlaying bool
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
		session: dg,
		cache:   NewLru(100),
		mutex:   sync.RWMutex{},
		wg:      sync.WaitGroup{},
		queue:   make([]string, 0),
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
			Log.Info("Message received")

			go func() {
				worker(b, m.Content, g.ID, vs.ChannelID)
			}()

			Log.Info("Successfully played")

			return
		}
	}
}

func (b *Bot) CreateTTS(message string) string {
	result := b.cache.Get(message)

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

	b.cache.Set(message, fileName)

	return fileName
}

func (b *Bot) PlayTTS(message, gID, vsChannelID string) {
	fmt.Println("isPlaying", b.isPlaying)
	if b.isPlaying == true {
		return
	}

	b.isPlaying = true

	var filePath string

	fmt.Println("1")
	b.wg.Add(1)
	go func() {
		filePath = b.CreateTTS(message)
		b.wg.Done()
	}()

	fmt.Println("2")
	b.wg.Wait()
	fmt.Println("3")

	if filePath == "" {
		Log.Warn("TTS sound doesn't created")
		return
	}

	fmt.Println("4")
	b.wg.Add(1)
	go func() {
		//sound := NewSound(filePath)
		b.Load(filePath, b.session, gID, vsChannelID)

		b.wg.Done()
	}()
	fmt.Println("5")
	b.wg.Wait()
	fmt.Println("6")

	b.QueueRemoveFisrt()
	b.isPlaying = false
	b.TimeoutCreate()
}

func (b *Bot) Stop() {

}

func (b *Bot) TimeoutCreate() {
	//timer := time.NewTimer(10 * time.Second)

	if len(b.queue) != 0 {
		Log.Info("Get sound from queue")
		b.PlayTTS(b.GetSound(), b.guildID, b.channelID)
	} else {
		Log.Info("Disconnect timer activated")
		//<-timer.C
	}
	Log.Info("Bot was inactive for 10 seconds")

	//b.voiceChannel.Disconnect()
	//b.voiceChannel = nil
}

func (b *Bot) globalTicker() {
	timer := time.NewTimer(900 * time.Second)
	go func() {
		<-timer.C
		b.voiceChannel.Disconnect()
	}()
}
