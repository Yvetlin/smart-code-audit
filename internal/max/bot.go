package max

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func RunIDBot(ctx context.Context, token string) error {
	api := NewAPI(token)
	var marker *int64

	log.Println("MAX бот запущен. Команды: /myid, /start, /id")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		updates, nextMarker, err := api.PollUpdates(marker, 30)
		if err != nil {
			log.Printf("poll error: %v", err)
			continue
		}
		if nextMarker != nil {
			marker = nextMarker
		}

		for _, raw := range updates {
			if err := handleUpdate(api, raw); err != nil {
				log.Printf("handle error: %v", err)
			}
		}
	}
}

func handleUpdate(api *API, raw json.RawMessage) error {
	var base struct {
		UpdateType string `json:"update_type"`
	}
	if err := json.Unmarshal(raw, &base); err != nil {
		return err
	}

	switch base.UpdateType {
	case "bot_started":
		var update struct {
			ChatID int64 `json:"chat_id"`
			User   struct {
				UserID int64 `json:"user_id"`
			} `json:"user"`
		}
		if err := json.Unmarshal(raw, &update); err != nil {
			return err
		}
		return api.SendToUser(update.User.UserID, formatIDMessage(update.User.UserID, update.ChatID))

	case "message_created":
		var update struct {
			Message struct {
				Sender *struct {
					UserID int64 `json:"user_id"`
					IsBot  bool  `json:"is_bot"`
				} `json:"sender"`
				Recipient struct {
					ChatID *int64 `json:"chat_id"`
				} `json:"recipient"`
				Body struct {
					Text *string `json:"text"`
				} `json:"body"`
			} `json:"message"`
		}
		if err := json.Unmarshal(raw, &update); err != nil {
			return err
		}
		if update.Message.Sender == nil || update.Message.Sender.IsBot {
			return nil
		}
		if update.Message.Body.Text == nil || !isIDCommand(*update.Message.Body.Text) {
			return nil
		}

		chatID := int64(0)
		if update.Message.Recipient.ChatID != nil {
			chatID = *update.Message.Recipient.ChatID
		}
		return api.SendToUser(update.Message.Sender.UserID, formatIDMessage(update.Message.Sender.UserID, chatID))
	}

	return nil
}

func isIDCommand(text string) bool {
	cmd := strings.Fields(strings.TrimSpace(text))[0]
	cmd = strings.Split(cmd, "@")[0]
	switch strings.ToLower(cmd) {
	case "/myid", "/start", "/id":
		return true
	default:
		return false
	}
}

func formatIDMessage(userID, chatID int64) string {
	msg := fmt.Sprintf(
		"Ваш **user_id**: `%d`\n\nДобавьте в secrets CI:\n`MAX_USER_ID=%d`",
		userID, userID,
	)
	if chatID != 0 {
		msg += fmt.Sprintf("\n\nДиалог (chat_id): `%d`", chatID)
	}
	return msg
}
