package main

import (
	"github.com/bwmarrin/discordgo"
	"local.packages/dca"
)

type Now struct {
	Author *discordgo.User
	Text   string
	Audio  *dca.EncodeSession
	IsLast bool
}
