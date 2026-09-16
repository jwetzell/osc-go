package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jwetzell/osc-go"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v3"
)

func main() {

	cmd := &cli.Command{
		Name:  "theatrechat",
		Usage: "receive theatre chat messages via OSC",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "ip",
				Usage: "ip to receive OSC messages on",
				Value: "0.0.0.0",
			},
			&cli.Int32Flag{
				Name:  "port",
				Usage: "port to receive OSC messages on",
				Value: 27900,
			},
			&cli.StringSliceFlag{
				Name:  "channel",
				Usage: "Startup channels",
				Value: []string{},
			},
			&cli.StringFlag{
				Name:  "username",
				Usage: "username to use for sending messages",
				Value: "User",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			ip := cmd.String("ip")
			port := cmd.Int32("port")
			startupChannels := cmd.StringSlice("channel")
			username := cmd.String("username")

			netAddress := fmt.Sprintf("%s:%d", ip, port)

			outConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
				IP:   net.ParseIP("255.255.255.255"),
				Port: int(port),
			})

			if err != nil {
				return err
			}

			chatApp := &ChatApp{
				username:        username,
				selectedChannel: "",
				app:             tview.NewApplication(),
				chatView:        tview.NewFlex(),
				channels:        []*Channel{},
				messageInput:    tview.NewInputField(),
				channelList:     tview.NewList(),
				mainFrame:       tview.NewFlex(),
				outConn:         outConn,
			}

			chatApp.initViews()

			if len(startupChannels) > 0 {
				for _, ch := range startupChannels {
					chatApp.addChannel(ch)
				}
			}

			//TODO(jwetzell): actually handle UDP loop lifecycle?
			go listenUDP(netAddress, chatApp.handleChat, chatApp.handleFlash)

			chatApp.app.EnableMouse(true)

			err = chatApp.app.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

type ChatApp struct {
	username        string
	selectedChannel string
	app             *tview.Application
	chatView        *tview.Flex
	channelList     *tview.List
	mainFrame       *tview.Flex
	messageInput    *tview.InputField
	channels        []*Channel
	outConn         *net.UDPConn
}

func (a *ChatApp) handleChat(chat Chat) {
	a.app.QueueUpdateDraw(func() {
		for _, c := range a.channels {
			if c.Name == chat.Channel {
				c.AddChat(chat)
				return
			}
		}
	})
}

func (a *ChatApp) handleFlash(channel string) {
	for _, c := range a.channels {
		if c.Name == channel && a.selectedChannel == channel {
			//TODO(jwetzell): flash channel
			return
		}
	}
}

func (a *ChatApp) initViews() {
	a.channelList.SetTitle("Channels")
	a.channelList.SetBorder(true)
	a.chatView.SetBorder(true)
	a.chatView.SetTitle("Messages (Select Channel)")
	a.chatView.SetDirection(tview.FlexRow)

	a.mainFrame.AddItem(a.channelList, 0, 1, false)
	a.mainFrame.AddItem(a.chatView, 0, 3, true)

	a.app.SetRoot(a.mainFrame, true).SetFocus(a.channelList)

	a.messageInput.SetPlaceholder("Message")
	a.messageInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter && a.selectedChannel != "" {
			a.sendChat()
		}
	})
}

func (a *ChatApp) channelListSelected(channel *Channel) {
	a.selectedChannel = channel.Name
	a.chatView.Clear()
	a.chatView.SetTitle(fmt.Sprintf("Messages (%s)", channel.Name))
	a.chatView.AddItem(channel.MessagesView, 0, 1, false)
	a.chatView.AddItem(a.messageInput, 1, 0, true)
	a.app.SetFocus(a.messageInput)
}

func (a *ChatApp) addChannel(name string) {
	channel := &Channel{
		Name:         name,
		MessagesView: tview.NewTextView().SetDynamicColors(true),
	}
	a.channels = append(a.channels, channel)
	a.channelList.AddItem(channel.Name, "", 0, func() {
		a.channelListSelected(channel)
	})
}

func (a *ChatApp) sendChat() {
	msg := a.messageInput.GetText()
	if msg != "" {
		chatMessage := osc.Message{
			Address: fmt.Sprintf("/theatrechat/message/%s", a.selectedChannel),
			Args: []osc.Arg{
				osc.StringArg(a.username),
				osc.StringArg(msg),
			},
		}
		bytes, err := chatMessage.ToBytes()
		if err == nil {
			a.outConn.Write(bytes)
		}
		a.messageInput.SetText("")
	}
}

func listenUDP(netAddress string, handleChat func(chat Chat), handleFlash func(channel string)) {

	laddr, err := net.ResolveUDPAddr("udp4", netAddress)
	if err != nil {
		return
	}

	conn, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return
	}

	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		bytesRead, _, err := conn.ReadFromUDP(buffer)

		if err != nil {
			panic(err)
		}

		oscPacket, _, err := osc.PacketFromBytes(buffer[0:bytesRead])

		if err != nil {
			panic(err)
		}

		msg, ok := oscPacket.(*osc.Message)
		if ok {
			if strings.HasPrefix(msg.Address, "/theatrechat/message/") {
				channel := strings.TrimPrefix(msg.Address, "/theatrechat/message/")
				if len(msg.Args) == 2 {
					user := msg.Args[0].Value
					text := msg.Args[1].Value
					handleChat(Chat{
						Channel: channel,
						User:    fmt.Sprintf("%v", user),
						Text:    fmt.Sprintf("%v", text),
					})
				}
			} else if strings.HasPrefix(msg.Address, "/theatrechat/flash/") {
				channel := strings.TrimPrefix(msg.Address, "/theatrechat/flash/")
				handleFlash(channel)
			}
		}
	}
}
