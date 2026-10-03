package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/cristiangonsevi/learning-go/portbridge/internal/cli"
)

var availableProfiles []string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cli.Execute(ctx)
}
