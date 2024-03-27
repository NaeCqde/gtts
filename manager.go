package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"local.packages/dca"
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
	Stoped  bool
}

func (m *VoiceManager) Play(texts []string, author *discordgo.User) {
	last := texts[len(texts)-1]

	for _, text := range texts[:len(texts)-1] {
		fmt.Printf("%v#%v %v\n", author.Username, author.Discriminator, text)

		m.Queue.Add(text, author, false)
	}

	fmt.Printf("%v#%v %v\n", author.Username, author.Discriminator, last)

	m.Queue.Add(last, author, true)

	m.work()
}

func (m *VoiceManager) Skip() {
	m.Queue.Delete(len(m.Queue.List))

	if m.Stream != nil {
		m.Stream.Stop()
		m.Stream.Cleanup()
	}
}

func (m *VoiceManager) Leave() {
	m.Stoped = true

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
	if !m.Queue.IsWorking {
		m.Queue.IsWorking = true

		var isLast bool

		for {
			if m.Stoped {
				return
			}

			if len(m.Queue.List) != 0 {
				now := m.Queue.List[0]

				isLast = now.IsLast
				m.Stream = now.Audio

				done := make(chan error)

				dca.NewStream(m.Stream, m.Conn, done)

				<-done

				m.Queue.Delete(1)

				m.Stream.Cleanup()
			}

			if isLast {
				time.Sleep(time.Duration(time.Millisecond * 10))
			}
		}
	}
}

func NewManager(conn *discordgo.VoiceConnection, guildID string, textChannel *discordgo.Channel, voiceChannel *discordgo.Channel) (*VoiceManager, error) {
	instance := &VoiceManager{Conn: conn, Channel: ChannelProps{Text: textChannel, Voice: voiceChannel}, GuildID: guildID, Queue: &Queue{IsWorking: false}}

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
