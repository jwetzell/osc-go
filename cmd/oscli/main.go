package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

var version = "local"

func main() {

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{})))
	cmd := &cli.Command{
		Name:    "oscli",
		Usage:   "basic OSC utilities",
		Version: version,
		Commands: []*cli.Command{
			newSendCommand(),
			newReceiveCommand(),
			newEncodeCommand(),
			newDecodeCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Default().Error("Error running command", "err", err)
	}
}
