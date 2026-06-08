package main

import (
	"fmt"
	"log"
	"os"

	"smart-code-audit/internal/max"
)

func main() {
	token := os.Getenv("MAX_BOT_TOKEN")
	if token == "" {
		log.Fatal("MAX_BOT_TOKEN is not set")
	}

	users, err := max.DiscoverUsers(token)
	if err != nil {
		log.Fatalf("failed to fetch updates: %v", err)
	}

	if len(users) == 0 {
		fmt.Println("Пользователи не найдены.")
		fmt.Println()
		fmt.Println("Что сделать:")
		fmt.Println("1. Запустите бота:  MAX_BOT_TOKEN=... go run ./cmd/max-bot")
		fmt.Println("2. Напишите боту /myid в MAX")
		fmt.Println("   или нажмите «Начать» — бот пришлёт ID сам")
		fmt.Println("3. Альтернатива: запустите эту команду снова после сообщения боту")
		os.Exit(1)
	}

	fmt.Println("Найденные пользователи (для личных уведомлений):")
	fmt.Println()
	for _, u := range users {
		fmt.Printf("  %s", u.Name)
		if u.Username != "" {
			fmt.Printf(" (@%s)", u.Username)
		}
		fmt.Printf("\n    user_id:  %d  ← положите в MAX_USER_ID\n", u.UserID)
		if u.ChatID != 0 {
			fmt.Printf("    chat_id:  %d  (диалог с ботом)\n", u.ChatID)
		}
		fmt.Printf("    событие:  %s\n\n", u.Event)
	}

	fmt.Println("Несколько получателей: MAX_USER_ID=123,456,789")
}
