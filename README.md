# gtts

日本語 | [English](README.en.md)

Discord のテキストチャンネルのメッセージを、ボイスチャンネルで読み上げる Go 製のボットです(discordgo + disgolf)。音声は Google 翻訳の TTS か、AITalk のデモ用 Web API(VOICEROID 系の声)で作ります。

## 使い方

1. リリースから実行ファイルをダウンロードします(Windows / macOS / Linux)。ffmpeg が必要です。
2. 設定の雛形 `example.env` を `.env` にして、Discord ボットのトークンなどを書きます。読み上げのブロックリストと声の設定の雛形は `example.blocklist.json` と `example.voice_settings.json` です。
3. 実行ファイルを起動します。

ビルドする場合は `go build` です。

## 仕様

動作の決まりは [docs/SPEC.md](docs/SPEC.md) にあります。
