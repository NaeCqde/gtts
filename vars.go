package main

import (
	"os"

	gojson "github.com/goccy/go-json"
	"github.com/joho/godotenv"
)

var (
	VoiceSettings = VoiceSettingMap{}
	BlackList     []string
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

func LoadBlackLists() error {
	body, err := os.ReadFile("blacklists.json")

	if err != nil {
		return err
	}

	return gojson.Unmarshal(body, &BlackList)
}

func SaveBlackLists() error {
	data, err := gojson.Marshal(BlackList)

	if err != nil {
		return err
	}

	return os.WriteFile("blacklists.json", data, os.ModePerm)
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
