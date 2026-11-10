package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/NaeCqde/dca"
)

var (
	SilenceDetectPattern = regexp.MustCompile(`(?m) silence_end: (\d+(?:\.\d+)?) \|`)
)

func GetAudio(text string, voiceType VoiceType) (*dca.EncodeSession, error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Channels = 1
	options.Bitrate = 64
	options.Application = dca.AudioApplicationLowDelay

	textLen := utf8.RuneCountInString(text)

	var audioURL string

	options.AudioFilter = "volume=-20dB;"

	if voiceType == 0 {
		audioURL = fmt.Sprintf(
			"https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&client=tw-ob&tl=%s&q=%s&textlen=%d",
			"ja",
			url.QueryEscape(text),
			textLen,
		)

		if textLen > 100 {
			options.AudioFilter += "atempo=2.4"
		} else if textLen > 50 {
			options.AudioFilter += "atempo=1.8"
		} else {
			options.AudioFilter += "atempo=1.2"
		}
	} else {
		speed := 1.2
		if textLen > 90 {
			speed = 2.4
		} else if textLen > 50 {
			speed = 1.8
		}

		data := url.Values{
			"speaker_id": {strconv.Itoa(voiceType)},
			"text":       {text},
			"ext":        {"ogg"},
			"volume":     {"1.2"},
			"speed":      {strconv.FormatFloat(speed, 'f', -1, 64)},
			"pitch":      {"1.0"},
			"range":      {"1.0"},
			"callback":   {"callback"},
		}

		resp, err := http.PostForm("https://cloud.ai-j.jp/demo/aitalk2webapi_nop.php", data)

		if err != nil {
			return &dca.EncodeSession{}, err
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

		if err != nil {
			return &dca.EncodeSession{}, err
		}

		audioURL = "https:" + strings.Replace(strings.Replace(strings.ReplaceAll(string(body), "\\", ""), `callback({"url":"`, "", 1), `"})`, "", 1)

		if u, perr := url.Parse(audioURL); perr != nil || u.Scheme != "https" || u.Host == "" || strings.HasPrefix(audioURL, "https:-") {
			return &dca.EncodeSession{}, fmt.Errorf("音声合成に失敗")
		}

		if strings.Contains(audioURL, "合成失敗") || audioURL == "https:" {
			return &dca.EncodeSession{}, fmt.Errorf("音声合成に失敗")
		}

		ffmpeg := exec.Command("ffmpeg", "-protocol_whitelist", "https,tls,tcp", "-i", audioURL, "-af", "silencedetect=n=-90dB:d=0.7", "-f", "null", "-")

		stderr, err := ffmpeg.CombinedOutput()

		if err != nil {
			return &dca.EncodeSession{}, err
		}

		result := SilenceDetectPattern.FindStringSubmatch(string(stderr))

		if len(result) < 2 {
			return &dca.EncodeSession{}, fmt.Errorf("無音検出の結果の取得に失敗しました")
		}

		startTime, err := strconv.ParseFloat(result[1], 64)

		if err != nil {
			return &dca.EncodeSession{}, fmt.Errorf("無音検出の結果の取得に失敗しました")
		}

		switch speed {
		case speed:
			// default
			options.StartTime = startTime
		case 2.4:
			options.StartTime = startTime / 2.4
		case 1.8:
			options.StartTime = startTime / 1.8
		}
	}

	return dca.EncodeFile(audioURL, options)
}
