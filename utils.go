package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	gojson "github.com/goccy/go-json"
)

func GetInvite(inviteCode string) (string, error) {
	resp, err := http.Get("https://discord-invite.naecqde.workers.dev/" + url.PathEscape(inviteCode))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return "", errors.New("discord ratelimited")
	}

		if resp.StatusCode != http.StatusOK {
		return "", errors.New("discord error")
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var data map[string]any

	err = gojson.Unmarshal(body, &data)

	if err != nil {
		return "", err
	}

	guild, ok := data["guild"].(map[string]any)
	if !ok {
		return "", errors.New("unexpected invite response")
	}

	name, ok := guild["name"].(string)
	if !ok {
		return "", errors.New("unexpected invite response")
	}

	return name, nil
}
