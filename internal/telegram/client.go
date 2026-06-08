package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	token  string
	chatID string
}

func New(tokenEnv, chatEnv string) *Client {
	return &Client{
		token:  os.Getenv(tokenEnv),
		chatID: os.Getenv(chatEnv),
	}
}

func (c *Client) Configured() bool {
	return c.token != "" && c.chatID != ""
}

func (c *Client) Send(text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)

	data := url.Values{}
	data.Set("chat_id", c.chatID)
	data.Set("text", text)
	data.Set("parse_mode", "Markdown")

	_, err := http.PostForm(apiURL, data)
	return err
}
