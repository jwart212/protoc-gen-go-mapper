package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/iktdev-boilerplate/go/protoc-gen-go-mapper/examples/advanced/internal/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer cancel()

	cmd, err := cmd.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("APP RUNNING")
	if err := cmd.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
