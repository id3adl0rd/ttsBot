package main

import (
	"time"
)

type Guild struct {
	actions *Actions
	media   chan *Media
	time    time.Time
	name    string
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
}

func (g *Guild) UpdateTime() {
	g.time = time.Now().Add(3 * time.Second)
}

func (g *Guild) ShouldBeDeleted() bool {
	return time.Now().After(g.time) || time.Now().Equal(g.time)
}

type Media struct {
	message string
	path    string
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
