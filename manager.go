package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/NaeCqde/dca"
	"github.com/bwmarrin/discordgo"
)

var Instances = []*VoiceManager{}

type ChannelProps struct {
	Text  *discordgo.Channel
	Voice *discordgo.Channel
}

type VoiceManager struct {
	Conn    *discordgo.VoiceConnection
	Channel ChannelProps
	GuildID string
	Queue   *Queue
	Stream  *dca.EncodeSession
	Stopped bool
}

func (m *VoiceManager) Play(texts []string, author *discordgo.User) {
	last := texts[len(texts)-1]

	for _, text := range texts[:len(texts)-1] {
		fmt.Printf("%v#%v %v\n", author.Username, author.Discriminator, text)

		m.Queue.Add(text, author, false)
	}

	fmt.Printf("%v#%v %v\n", author.Username, author.Discriminator, last)

	m.Queue.Add(last, author, true)

	if !m.Queue.IsWorking {
		go m.work()
	}
}

func (m *VoiceManager) Skip() {
	m.Queue.Delete(len(m.Queue.List))

	if m.Stream != nil {
		m.Stream.Cleanup()
	}
}

func (m *VoiceManager) Leave() {
	m.Stopped = true

	m.Conn.Disconnect()

	new := []*VoiceManager{}

	for _, vm := range Instances {
		if m.GuildID != vm.GuildID {
			new = append(new, vm)
		}
	}

	Instances = new
}

func (m *VoiceManager) work() {
	m.Queue.IsWorking = true

	for len(m.Queue.List) != 0 {
		if m.Stopped {
			m.Queue.IsWorking = false
			return
		}

		now := m.Queue.List[0]

		m.Stream = now.Audio
		done := make(chan error)

		dca.NewStream(m.Stream, m.Conn, done)

		err := <-done
		if err != nil {
			slog.Error(err.Error())
		}

		m.Queue.Delete(1)
		m.Stream.Cleanup()

		if now.IsPartLast {
			time.Sleep(time.Duration(time.Millisecond * 10))
		}
	}

	m.Queue.IsWorking = false
}

func NewManager(
	conn *discordgo.VoiceConnection,
	guildID string, textChannel *discordgo.Channel,
	voiceChannel *discordgo.Channel,
) (*VoiceManager, error) {
	instance := &VoiceManager{
		Conn: conn,
		Channel: ChannelProps{
			Text:  textChannel,
			Voice: voiceChannel,
		},
		GuildID: guildID,
		Queue:   &Queue{IsWorking: false},
	}

	Instances = append(Instances, instance)

	return instance, nil
}

func GetManager(guildID string) (*VoiceManager, error) {
	for _, instance := range Instances {
		if guildID == instance.GuildID {
			return instance, nil
		}
	}

	return &VoiceManager{}, fmt.Errorf("failed to find Voice Manager")
}
