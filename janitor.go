package main

import (
	"fmt"
	"time"
)

const (
	leaveTime = 5 * time.Second
)

func cleanConnectionWithEmptyQueue(b *Bot) {
	for {
		cleanAllConnection(b)
		time.Sleep(leaveTime)
	}
}

func cleanAllConnection(b *Bot) {
	for _, connection := range b.session.VoiceConnections {
		guild := b.GetGuild(connection.GuildID)
		//fmt.Println("t1", guild.IsStreaming(), guild.time)
		//&& guild.ShouldBeDeleted()
		if !guild.IsStreaming() {
			fmt.Println("lolzxcasd", guild.ShouldBeDeleted())
			if guild.ShouldBeDeleted() {
				connection.Disconnect()
			}
		}
	}
}
