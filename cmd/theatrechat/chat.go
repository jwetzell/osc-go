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
