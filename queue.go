package main

import (
	"github.com/bwmarrin/discordgo"
	//"time"
)

type Queue struct {
	List []Now
	//lastDate  time.Time
	IsWorking bool
}

func (q *Queue) Add(text string, author *discordgo.User, isLast bool) {
	var voiceType VoiceType = 1209

	if setting, ok := VoiceSettings[author.ID]; ok {
		voiceType = setting.VoiceType
	}

	audio, err := GetAudio(text, voiceType)

	if err != nil {
		return
	}

	now := Now{
		Text:   text,
		Author: author,
		Audio:  audio,
		IsLast: isLast,
	}

	q.List = append(q.List, now)
}

func (q *Queue) Delete(amount int) {
	defer func() {
		recover()
	}()

	q.List = q.List[amount:]
}
