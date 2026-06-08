package max

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DiscoveredUser struct {
	UserID   int64
	ChatID   int64
	Name     string
	Username string
	Event    string
}

func DiscoverUsers(token string) ([]DiscoveredUser, error) {
	req, err := http.NewRequest(http.MethodGet, apiBaseURL+"/updates?limit=100", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("max api returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Updates []json.RawMessage `json:"updates"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	seen := make(map[int64]struct{})
	var users []DiscoveredUser

	for _, raw := range payload.Updates {
		var base struct {
			UpdateType string `json:"update_type"`
		}
		if err := json.Unmarshal(raw, &base); err != nil {
			continue
		}

		switch base.UpdateType {
		case "bot_started":
			var update struct {
				ChatID int64 `json:"chat_id"`
				User   struct {
					UserID    int64   `json:"user_id"`
					FirstName string  `json:"first_name"`
					LastName  *string `json:"last_name"`
					Username  *string `json:"username"`
				} `json:"user"`
			}
			if err := json.Unmarshal(raw, &update); err != nil {
				continue
			}
			if u, ok := addUser(seen, update.User.UserID, update.ChatID, update.User.FirstName, update.User.LastName, update.User.Username, "bot_started"); ok {
				users = append(users, u)
			}

		case "message_created":
			var update struct {
				Message struct {
					Sender *struct {
						UserID    int64   `json:"user_id"`
						FirstName string  `json:"first_name"`
						LastName  *string `json:"last_name"`
						Username  *string `json:"username"`
					} `json:"sender"`
					Recipient struct {
						ChatID *int64 `json:"chat_id"`
					} `json:"recipient"`
				} `json:"message"`
			}
			if err := json.Unmarshal(raw, &update); err != nil || update.Message.Sender == nil {
				continue
			}
			sender := update.Message.Sender
			chatID := int64(0)
			if update.Message.Recipient.ChatID != nil {
				chatID = *update.Message.Recipient.ChatID
			}
			if u, ok := addUser(seen, sender.UserID, chatID, sender.FirstName, sender.LastName, sender.Username, "message_created"); ok {
				users = append(users, u)
			}
		}
	}

	return users, nil
}

func addUser(seen map[int64]struct{}, userID, chatID int64, firstName string, lastName, username *string, event string) (DiscoveredUser, bool) {
	if userID == 0 {
		return DiscoveredUser{}, false
	}
	if _, exists := seen[userID]; exists {
		return DiscoveredUser{}, false
	}
	seen[userID] = struct{}{}

	name := firstName
	if lastName != nil && *lastName != "" {
		name += " " + *lastName
	}

	u := DiscoveredUser{
		UserID: userID,
		ChatID: chatID,
		Name:   name,
		Event:  event,
	}
	if username != nil {
		u.Username = *username
	}
	return u, true
}
