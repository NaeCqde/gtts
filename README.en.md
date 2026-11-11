# gtts

[日本語](README.md) | English

A bot written in Go (discordgo + disgolf) that reads the messages of a Discord text channel aloud in a voice channel. The audio is made with Google Translate TTS or the demo web API of AITalk (VOICEROID-type voices).

## Usage

1. Download the executable from the releases (Windows / macOS / Linux). ffmpeg is required.
2. Copy the config template `example.env` to `.env` and write the Discord bot token and so on. The templates of the read-aloud block list and the voice settings are `example.blocklist.json` and `example.voice_settings.json`.
3. Start the executable.

To build it yourself, use `go build`.

## Specification

The rules of the behavior are in [docs/SPEC.md](docs/SPEC.md) (in Japanese).
