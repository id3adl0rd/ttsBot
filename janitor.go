package main

import (
	"time"
)

func clearConnection(bot *Bot) {
	for {
		for _, connection := range bot.session.VoiceConnections {
			g := bot.GetGuild(connection.ChannelID)
			if g.QueueSize() == 0 && g.ShouldBeDeleted() {
				if g.ShouldBeDeleted() {
					connection.Disconnect()
				}
			}
		}

		time.Sleep(time.Duration(Enviroment.Misc.DisconnectTimer * int64(time.Second)))
	}
}
