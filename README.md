[![codecov](https://codecov.io/gh/jwetzell/osc-go/branch/main/graph/badge.svg?token=DW3CCZELGI)](https://codecov.io/gh/jwetzell/osc-go)

A mostly complete OSC implementation and collection of command line OSC utilities written in Go. Mainly an exercise in learning Go.

## Utilities

### sendosc
```
NAME:
   sendosc - send OSC messages via UDP or TCP

USAGE:
   sendosc [global options]

GLOBAL OPTIONS:
   --host string                    host to send OSC message to
   --port int                       port to send OSC message to
   --protocol string                protocol to use to send (tcp or udp) (default: "udp")
   --address string                 OSC address
   --arg string [ --arg string ]    OSC args
   --type string [ --type string ]  OSC types
   --slip                           whether to slip encode the OSC Message bytes
   --help, -h                       show help
```
### makeosc
```
NAME:
   makeosc - make osc bytes

USAGE:
   makeosc [global options]

GLOBAL OPTIONS:
   --address string                 OSC address
   --arg string [ --arg string ]    OSC args
   --type string [ --type string ]  OSC types
   --slip                           whether to slip encode the OSC Message bytes
   --help, -h                       show help
```
### receiveosc
```
NAME:
   receiveosc - receive OSC messages via UDP or TCP

USAGE:
   receiveosc [global options]

GLOBAL OPTIONS:
   --ip string        ip to receive OSC messages on (default: "0.0.0.0")
   --port int         port to receive OSC messages on (default: 8888)
   --protocol string  protocol to use to receive (tcp or udp) (default: "udp")
   --format string    format for messages to be output in ('json') (default: "json")
   --slip             whether to slip encode the OSC Message bytes
   --help, -h         show help
```
