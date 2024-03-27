package main

import (
	"errors"
	"io"
	"net/http"

	gojson "github.com/goccy/go-json"
)

func GetInvite(inviteCode string) (string, error) {
	resp, err := http.Get("https://discord-invite.naecqde.workers.dev/" + inviteCode)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("discord ratelimited or error")
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

	return data["guild"].(map[string]any)["name"].(string), nil
}
