package main

type VoiceType = int

var (
	VoiceTypeMap = map[VoiceType]string{
		0:    "Google TTS",
		203:  "VOICEROID+ 琴葉 葵",
		551:  "VOICEROID2 琴葉 葵",
		552:  "VOICEROID2 琴葉 茜",
		2006: "VOICEROID2 音街 ウナ",
		554:  "VOICEROID2 紲星 あかり",
		553:  "VOICEROID2 結月 ゆかり",
		1001: "VOICEROID+ 琴葉 茜",
		1201: "VOICEROID+ 鷹の爪 吉田くん",
		1202: "VOICEROID+ 東北 ずん子",
		1203: "VOICEROID+ 月読アイ",
		1204: "VOICEROID+ 月読ショウタ",
		1205: "VOICEROID+ 民安ともえ",
		1206: "VOICEROID+ 結月 ゆかり",
		1207: "VOICEROID+ 水奈瀬 コウ",
		1208: "VOICEROID+ 京町 セイカ",
		1209: "VOICEROID+ 東北 きりたん",
		1210: "VOICEROID2 桜乃 そら",
		1211: "VOICEROID2 東北 イタコ",
		1212: "VOICEROID2 標準語 ついなちゃん",
		1213: "VOICEROID2 関西弁 ついなちゃん",
		1214: "VOICEROID2 伊織 弓鶴",
	}
	VoiceTypeRow = []int{
		0,
		551,
		552,
		203,
		1001,
		1203,
		1204,
		1210,
		2006,
		1214,
		554,
		1211,
		553,
		1202,
		1207,
		1206,
		1208,
		1209,
		1205,
		1201,
		1212,
		1213,
	}
)

type VoiceSetting struct {
	VoiceType VoiceType `json:"voice_type"`
}

type VoiceSettingMap map[string]VoiceSetting
