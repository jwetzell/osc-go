package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	osc "github.com/jwetzell/osc-go"
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			ip := cmd.String("ip")
			port := cmd.Int32("port")

			netAddress := fmt.Sprintf("%s:%d", ip, port)

			outConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
				IP:   net.ParseIP("255.255.255.255"),
				Port: int(port),
			})
			if err != nil {
				return err
			}

			channels := []*Channel{}

			app := tview.NewApplication()
			mainFrame := tview.NewFlex()
			channelList := tview.NewList()
			chatView := tview.NewFlex()
			messageInput := tview.NewInputField()

			channelList.SetTitle("Channels")
			channelList.SetBorder(true)
			chatView.SetBorder(true)
			chatView.SetTitle("Messages (Select Channel)")
			chatView.SetDirection(tview.FlexRow)

			mainFrame.AddItem(channelList, 0, 1, false)
			mainFrame.AddItem(chatView, 0, 3, true)

			app.SetRoot(mainFrame, true).SetFocus(channelList)

			selectedChannel := ""

			messageInput.SetPlaceholder("Message")
			messageInput.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter && selectedChannel != "" {
					text := messageInput.GetText()
					if text != "" {
						// do something
						chatMessage := osc.Message{
							Address: fmt.Sprintf("/theatrechat/message/%s", selectedChannel),
							Args: []osc.Arg{
								osc.StringArg("Me"),
								osc.StringArg(text),
							},
						}
						bytes, err := chatMessage.ToBytes()
						if err == nil {
							outConn.Write(bytes)
						}
						messageInput.SetText("")
					}
				}
			})

			go listenUDP(netAddress, func(chat Chat) {
				app.QueueUpdateDraw(func() {

					var channel *Channel
					for _, c := range channels {
						if c.Name == chat.Channel {
							channel = c
							break
						}
					}
					if channel == nil {
						channel = &Channel{
							Name:         chat.Channel,
							MessagesView: tview.NewTextView().SetDynamicColors(true),
						}
						channels = append(channels, channel)
					}

					fmt.Fprintf(channel.MessagesView, "[%s] %s: %s\n", chat.Channel, chat.User, chat.Text)

					channelList.Clear()
					for _, channel := range channels {
						ch := channel
						channelList.AddItem(ch.Name, "", 0, func() {
							selectedChannel = ch.Name
							chatView.Clear()
							chatView.SetTitle(fmt.Sprintf("Messages (%s)", ch.Name))
							chatView.AddItem(ch.MessagesView, 0, 1, false)
							chatView.AddItem(messageInput, 1, 0, true)
							app.SetFocus(messageInput)
						})
					}

				})
			})

			err = app.Run()
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

type Channel struct {
	Name         string
	MessagesView *tview.TextView
}

type Chat struct {
	Channel string
	User    string
	Text    string
}

func listenUDP(netAddress string, handleChat func(chat Chat)) {

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
			}
		}
	}
}
