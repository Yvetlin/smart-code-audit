package notify

import (
	"fmt"
	"os"

	"smart-code-audit/internal/config"
	"smart-code-audit/internal/max"
	"smart-code-audit/internal/telegram"
)

type Notifier interface {
	Send(text string) error
	Channel() string
}

func New(cfg *config.Config) Notifier {
	maxClient := max.New(cfg.Max.BotTokenEnv, cfg.Max.UserIDEnv, cfg.Max.ChatIDEnv)
	if maxClient.Configured() {
		return &channelNotifier{name: "MAX", client: maxClient}
	}

	tgClient := telegram.New(cfg.Telegram.BotTokenEnv, cfg.Telegram.ChatIDEnv)
	if tgClient.Configured() {
		return &channelNotifier{name: "Telegram", client: tgClient}
	}

	return nil
}

func SkipReason(cfg *config.Config) string {
	maxToken := os.Getenv(cfg.Max.BotTokenEnv)
	maxUserID := os.Getenv(cfg.Max.UserIDEnv)
	maxChatID := os.Getenv(cfg.Max.ChatIDEnv)
	tgToken := os.Getenv(cfg.Telegram.BotTokenEnv)
	tgChatID := os.Getenv(cfg.Telegram.ChatIDEnv)

	if maxToken != "" && maxUserID == "" && maxChatID == "" {
		return fmt.Sprintf(
			"MAX_BOT_TOKEN задан, но нет %s (личка) или %s (чат)",
			cfg.Max.UserIDEnv, cfg.Max.ChatIDEnv,
		)
	}
	if tgToken != "" && tgChatID == "" {
		return fmt.Sprintf("TG_BOT_TOKEN задан, но нет %s", cfg.Telegram.ChatIDEnv)
	}
	return "задайте MAX_BOT_TOKEN+MAX_USER_ID или TG_BOT_TOKEN+TG_CHAT_ID"
}

type sender interface {
	Send(text string) error
}

type channelNotifier struct {
	name   string
	client sender
}

func (n *channelNotifier) Send(text string) error {
	return n.client.Send(text)
}

func (n *channelNotifier) Channel() string {
	return n.name
}
