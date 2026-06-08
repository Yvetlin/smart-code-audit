package max

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type API struct {
	token string
}

func NewAPI(token string) *API {
	return &API{token: token}
}

func (a *API) SendToUser(userID int64, text string) error {
	return a.send("user_id", strconv.FormatInt(userID, 10), text)
}

func (a *API) send(param, id, text string) error {
	apiURL := fmt.Sprintf("%s/messages?%s=%s", apiBaseURL, param, id)

	body, err := json.Marshal(map[string]string{
		"text":   text,
		"format": "markdown",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", a.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("max api returned status %d", resp.StatusCode)
	}

	return nil
}

func (a *API) PollUpdates(marker *int64, timeout int) ([]json.RawMessage, *int64, error) {
	url := fmt.Sprintf("%s/updates?timeout=%d&types=message_created,bot_started", apiBaseURL, timeout)
	if marker != nil {
		url += "&marker=" + strconv.FormatInt(*marker, 10)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", a.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, nil, fmt.Errorf("max api returned status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var payload struct {
		Updates []json.RawMessage `json:"updates"`
		Marker  *int64            `json:"marker"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, nil, err
	}

	return payload.Updates, payload.Marker, nil
}
