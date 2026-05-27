package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GabsMeloTI/go_adega/cmd"
	"github.com/GabsMeloTI/go_adega/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	container, err := config.NewContainer(config.Load())
	if err != nil {
		log.Fatalf("cannot start container: %v", err)
	}
	defer container.Close()

	cmd.StartAPI(ctx, container)
}
