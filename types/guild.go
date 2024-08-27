package types

import (
	"time"
)

type Guild struct {
	actions   *Actions
	media     chan *Media
	time      time.Time
	name      string
	isPlaying bool
}

func NewGuild(name string) *Guild {
	return &Guild{name: name}
}

func (g *Guild) Enqueue(media *Media) {
	g.media <- media
}

func (g *Guild) IsStreaming() bool {
	return g.media != nil
}

func (g *Guild) QueueSize() int {
	return len(g.media)
}

func (g *Guild) PrepareMediaChannel(size int) {
	g.media = make(chan *Media, size)
	g.actions = NewActions()
}

func (g *Guild) IsQueueFull() bool {
	return g.media != nil && len(g.media) == cap(g.media)
}

func (g *Guild) Stop() {
	close(g.media)
	g.media = nil
	g.actions = nil
	g.isPlaying = false
}

func (g *Guild) UpdateTime(t int64) {
	g.time = time.Now().Add(time.Duration(t * int64(time.Second)))
}

func (g *Guild) ShouldBeDeleted() bool {
	return time.Now().After(g.time) || time.Now().Equal(g.time)
}

func (g *Guild) IsPlaying() bool {
	return g.isPlaying
}

func (g *Guild) GetMedia() chan *Media {
	return g.media
}

func (g *Guild) SetIsPlaying(v bool) {
	g.isPlaying = v
}

type Media struct {
	message string
	path    string
}

func (m *Media) SetMessage(message string) {
	m.message = message
}

func (m *Media) SetPath(path string) {
	m.path = path
}

func (m *Media) GetMessage() string {
	return m.message
}

func (m *Media) GetPath() string {
	return m.path
}

func NewMedia(message string, path string) *Media {
	return &Media{message: message, path: path}
}

type Actions struct {
	stopChan chan bool
}

func NewActions() *Actions {
	return &Actions{
		stopChan: make(chan bool, 1),
	}
}

func (a *Actions) Stop() {
	a.stopChan <- true
}
