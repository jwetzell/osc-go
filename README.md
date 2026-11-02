[![codecov](https://codecov.io/gh/jwetzell/osc-go/branch/main/graph/badge.svg?token=DW3CCZELGI)](https://codecov.io/gh/jwetzell/osc-go)

A mostly complete OSC implementation and a simple utility CLI written in Go. Mainly an exercise in learning Go.

## `oscli`

### Usage
```
NAME:
   oscli - basic OSC utilities

USAGE:
   oscli [global options] [command [command options]]

COMMANDS:
   send, s     send OSC messages via UDP or TCP
   receive, r  receive OSC messages via UDP or TCP
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### `send`
```
NAME:
   oscli send - send OSC messages via UDP or TCP

USAGE:
   oscli send [options]

OPTIONS:
   --host string                    host to send OSC message to
   --port int                       port to send OSC message to
   --protocol string                protocol to use to send (tcp or udp) (default: "udp")
   --address string                 OSC address
   --arg string [ --arg string ]    OSC args
   --type string [ --type string ]  OSC types
   --slip                           whether to slip encode the OSC Message bytes
   --help, -h                       show help
```

### `receive`
```
NAME:
   oscli receive - receive OSC messages via UDP or TCP

USAGE:
   oscli receive [options]

OPTIONS:
   --ip string        ip to receive OSC messages on (default: "0.0.0.0")
   --port int         port to receive OSC messages on (default: 8888)
   --protocol string  protocol to use to receive (tcp or udp) (default: "udp")
   --format string    format for messages to be output in ('json') (default: "json")
   --slip             whether to slip encode the OSC Message bytes
   --help, -h         show help
```