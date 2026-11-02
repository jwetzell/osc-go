[![codecov](https://codecov.io/gh/jwetzell/osc-go/branch/main/graph/badge.svg?token=DW3CCZELGI)](https://codecov.io/gh/jwetzell/osc-go)

A mostly complete OSC implementation and collection of command line OSC utilities written in Go. Mainly an exercise in learning Go.

## Utilities

### sendosc
```
Usage of sendosc:
  -address string
        OSC address
  -arg value
        OSC args
  -host string
        host to send OSC message to
  -port int
        port to send OSC message to
  -protocol string
        protocol to use to send (tcp or udp) (default "udp")
  -slip
        whether to slip encode the OSC Message bytes
  -type value
        OSC types
```
### receiveosc
```
Usage of receiveosc:
  -format string
        format for messages to be output in ('json') (default "json")
  -ip string
        ip to receive OSC messages on (default "0.0.0.0")
  -port int
        port to receive OSC messages on (default 8888)
  -protocol string
        protocol to use to receive (tcp or udp) (default "udp")
  -slip
        whether to slip encode the OSC Message bytes
```
