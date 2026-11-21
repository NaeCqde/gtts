package main

import (
	"os"
	"path/filepath"

	gojson "github.com/goccy/go-json"
	"github.com/joho/godotenv"
)

var (
	VoiceSettings = map[string]VoiceSetting{}
	BlockList     []string
)

// settingsPath is the path of a settings file: the one in the current directory; only when that does not exist, the one next
// to the executable (the one that is found is the only one that is read). With neither, the current directory (where it is saved).
func settingsPath(name string) string {
	if _, err := os.Stat(name); err == nil {
		return name
	}

	if exe, err := os.Executable(); err == nil {
		if real, err := filepath.EvalSymlinks(exe); err == nil {
			exe = real
		}

		beside := filepath.Join(filepath.Dir(exe), name)

		if _, err := os.Stat(beside); err == nil {
			return beside
		}
	}

	return name
}

func LoadEnv() error {
	return godotenv.Load(settingsPath(".env"))
}

func LoadVoiceSettings() error {
	body, err := os.ReadFile(settingsPath("voice_settings.json"))

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

	return os.WriteFile(settingsPath("voice_settings.json"), data, 0o600)
}

func LoadBlockList() error {
	body, err := os.ReadFile(settingsPath("blocklist.json"))

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

	return os.WriteFile("blocklists.json", data, 0o600)
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
