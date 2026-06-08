package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"smart-code-audit/internal/max"
)

func main() {
	token := os.Getenv("MAX_BOT_TOKEN")
	if token == "" {
		log.Fatal("MAX_BOT_TOKEN is not set")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Запуск бота для привязки пользователей...")
	log.Println("Попросите коллег написать боту /myid в MAX")
	if err := max.RunIDBot(ctx, token); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}
