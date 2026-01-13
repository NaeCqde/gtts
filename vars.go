package main

import (
	"os"

	gojson "github.com/goccy/go-json"
	"github.com/joho/godotenv"
)

var (
	VoiceSettings = map[string]VoiceSetting{}
	BlockList     []string
)

func LoadEnv() error {
	return godotenv.Load()
}

func LoadVoiceSettings() error {
	body, err := os.ReadFile("voice_settings.json")

	if err != nil {
		return err
	}

	return gojson.Unmarshal(body, &VoiceSettings)
}

func SaveVoiceSettings() error {
	data, err := gojson.Marshal(VoiceSettings)

	if err != nil {
		return err
	}

	return os.WriteFile("voice_settings.json", data, os.ModePerm)
}

func LoadBlockList() error {
	body, err := os.ReadFile("blocklist.json")

	if err != nil {
		return err
	}

	return gojson.Unmarshal(body, &BlockList)
}

func SaveBlockLists() error {
	data, err := gojson.Marshal(BlockList)

	if err != nil {
		return err
	}

	return os.WriteFile("blocklists.json", data, os.ModePerm)
}

func GetToken() string {
	return os.Getenv("TOKEN")
}

func GetOwnerID() string {
	return os.Getenv("OWNER_ID")
}

func GetPrefix() string {
	return os.Getenv("PREFIX")
}
