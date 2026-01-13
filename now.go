package main

import (
	"github.com/NaeCqde/dca"
	"github.com/bwmarrin/discordgo"
)

type Now struct {
	Author     *discordgo.User
	Text       string
	Audio      *dca.EncodeSession
	IsPartLast bool
}
