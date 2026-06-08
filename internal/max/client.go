package max

import (
	"fmt"
	"os"
	"strings"
)

const apiBaseURL = "https://platform-api.max.ru"

type Client struct {
	token   string
	userIDs []string
	chatID  string
}

func New(tokenEnv, userIDEnv, chatEnv string) *Client {
	userIDs := parseIDs(os.Getenv(userIDEnv))
	return &Client{
		token:   os.Getenv(tokenEnv),
		userIDs: userIDs,
		chatID:  os.Getenv(chatEnv),
	}
}

func parseIDs(raw string) []string {
	parts := strings.Split(raw, ",")
	ids := make([]string, 0, len(parts))
	for _, part := range parts {
		id := strings.TrimSpace(part)
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func (c *Client) Configured() bool {
	if c.token == "" {
		return false
	}
	return len(c.userIDs) > 0 || c.chatID != ""
}

func (c *Client) Send(text string) error {
	if len(c.userIDs) > 0 {
		for _, userID := range c.userIDs {
			if err := c.sendTo("user_id", userID, text); err != nil {
				return fmt.Errorf("user %s: %w", userID, err)
			}
		}
		return nil
	}

	return c.sendTo("chat_id", c.chatID, text)
}

func (c *Client) sendTo(param, id, text string) error {
	return NewAPI(c.token).send(param, id, text)
}
