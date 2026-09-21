package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type Channel struct {
	Name         string
	MessagesView *tview.TextView
}

func (c *Channel) AddChat(chat Chat) {
	//TODO(jwetzell): message formatting?
	fmt.Fprintf(c.MessagesView, "[%s] %s: %s\n", chat.Channel, chat.User, chat.Text)
	c.MessagesView.ScrollToEnd()
}

type Chat struct {
	Channel string
	User    string
	Text    string
}

// structs from official TheatreChat session data
type channel struct {
	Id             string
	FriendlyName   string
	Colour         string
	FlashOnMessage bool
}

type session struct {
	Username     string
	Channels     []channel
	QuickSelects []string
	TextSize     int
	ColorTheme   int
	BroadcastIp  string
}
