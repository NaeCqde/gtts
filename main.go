package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"local.packages/disgolf"
)

func main() {
	slog.Info("Start")

	err := LoadEnv()

	if err != nil {
		panic(err)
	}

	err = LoadVoiceSettings()

	if err != nil {
		panic(err)
	}

	err = LoadBlackLists()

	if err != nil {
		panic(err)
	}

	bot, err := disgolf.New(GetToken())

	if err != nil {
		panic(err)
	}

	bot.Identify.Intents = discordgo.IntentsAll

	bot.Router.Register(&disgolf.Command{
		Name:               "join",
		Description:        "VCに参加して読み上げを開始します",
		MessageHandler:     disgolf.MessageHandlerFunc(HandleVCJoin),
		MessageMiddlewares: []disgolf.MessageHandler{disgolf.MessageHandlerFunc(HasOwnerMiddleware), disgolf.MessageHandlerFunc(IgnoreSelfMiddleware)},
	})

	bot.Router.Register(&disgolf.Command{
		Name:               "leave",
		Description:        "読み上げを終了してVCから退出します",
		MessageHandler:     disgolf.MessageHandlerFunc(HandleVCLeave),
		MessageMiddlewares: []disgolf.MessageHandler{disgolf.MessageHandlerFunc(HasOwnerMiddleware), disgolf.MessageHandlerFunc(IgnoreSelfMiddleware)},
	})

	bot.Router.Register(&disgolf.Command{
		Name:               "voice",
		Description:        "あなたのメッセとステータスを読み上げる音声を設定します",
		MessageHandler:     disgolf.MessageHandlerFunc(HandleVoiceSetting),
		MessageMiddlewares: []disgolf.MessageHandler{disgolf.MessageHandlerFunc(IgnoreSelfMiddleware), disgolf.MessageHandlerFunc(IgnoreBlackListedUser)},
	})

	bot.Router.Register(&disgolf.Command{
		Name:               "blacklist",
		Description:        "グローバルで読み上げさせないユーザーを管理します",
		MessageHandler:     disgolf.MessageHandlerFunc(HandleBlackLists),
		MessageMiddlewares: []disgolf.MessageHandler{disgolf.MessageHandlerFunc(HasOwnerMiddleware), disgolf.MessageHandlerFunc(IgnoreSelfMiddleware)},
	})

	bot.AddHandler(bot.Router.MakeMessageHandler(&disgolf.MessageHandlerConfig{
		Prefixes:      []string{GetPrefix()},
		MentionPrefix: false,
	}))

	bot.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		slog.Info("Bot is ready", "UserName", fmt.Sprintf("%v#%v", m.User.Username, m.User.Discriminator), "UserID", m.User.ID)
	})

	bot.AddHandler(HandleMessages)
	bot.AddHandler(HandleVoiceUpdates)

	if err = bot.Open(); err != nil {
		panic(err)
	}

	stopBot := make(chan os.Signal, 1)

	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-stopBot

	err = bot.Close()

	if err != nil {
		panic(err)
	}

	slog.Info("End")
}
