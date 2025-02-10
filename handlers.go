package main

import (
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"local.packages/disgolf"
)

var (
	EmojiPattern         = regexp.MustCompile(`(?m)<:(\w+):\d+>`)
	DiscordFilePattern   = regexp.MustCompile(`(?m)https?://(?:cdn|media).discord(?:app)?.(?:com|net)?/attachments/\d+/\d+/([\w\d\.-]+)(?:\?.+)?`)
	LinkPattern          = regexp.MustCompile(`(?mi)https?://(?:www\.)?([-a-z0-9@:%._\+~#=]{1,256}\.[a-z0-9()]{1,6}\b)(?:[-a-z0-9()!@:%_\+.~#?&//=]*)`)
	DiscordInvitePattern = regexp.MustCompile(`(?mi)(?:https?://)?(?:\w+\.)?discord\.(?:com/invite/|gg/)([a-z\d-]+)`)
	DuplicatedPattern    = regexp.MustCompile(`(?mi)([!！]{3,}|[\?？]{3,}|[wｗ]{3,}|[-ー]{3,}|[_＿]{3,}|[\^|＾]{3,})`)
	CodeBlockPattern     = regexp.MustCompile(`(?ms)\x60\x60\x60.*\x60\x60\x60`)
	ExcludePattern       = regexp.MustCompile(`(?m)\*`)
	BotPrefixesPattern   = regexp.MustCompile(`^\w{0,2}[!"#$%&'\*\+,\./;<=>?~]{1,2}\w+\s*(?:[\w[:punct:]]+\s*)*`)
)

func HasOwnerMiddleware(ctx *disgolf.MessageCtx) {
	fmt.Printf("%+v\n", ctx.Message)
	if GetOwnerID() == ctx.Message.Author.ID {
		ctx.Next()
	}
}

func IgnoreSelfMiddleware(ctx *disgolf.MessageCtx) {
	if ctx.Message.Author.ID != ctx.State.User.ID && !ctx.Message.Author.Bot {
		ctx.Next()
	}
}

func IgnoreBlackListedUser(ctx *disgolf.MessageCtx) {
	if !slices.Contains(BlockList, ctx.Message.Author.ID) {
		ctx.Next()
	}
}

func HandleVCJoin(ctx *disgolf.MessageCtx) {
	textChannel, err := ctx.Channel(ctx.Message.ChannelID)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	voiceManager, err := GetManager(ctx.Message.GuildID)

	voiceChannelID := "0"

	if len(ctx.Message.MentionChannels) != 0 {
		voiceChannelID = ctx.Message.MentionChannels[0].ID
	} else {
		for _, g := range ctx.State.Guilds {
			for _, v := range g.VoiceStates {
				if v.UserID == ctx.Message.Author.ID {
					voiceChannelID = v.ChannelID
				}
			}
		}
	}

	if err == nil && voiceManager.Channel.Voice.ID != voiceChannelID {
		voiceManager.Conn.ChangeChannel(voiceChannelID, false, true)
		return
	}

	voiceChannel, err := ctx.Channel(voiceChannelID)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	conn, err := ctx.ChannelVoiceJoin(ctx.Message.GuildID, voiceChannelID, false, false)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("参加", "ChannelID", voiceChannelID)

	NewManager(conn, ctx.Message.GuildID, textChannel, voiceChannel)
}

func HandleVCLeave(ctx *disgolf.MessageCtx) {
	voiceManager, err := GetManager(ctx.Message.GuildID)

	if err == nil {
		slog.Info("退出", "ChannelID", voiceManager.Channel.Voice.ID)

		voiceManager.Leave()
	}
}

func HandleVoiceSetting(ctx *disgolf.MessageCtx) {
	if len(ctx.Arguments) == 0 {
		text := "音声リスト\n\n"

		for _, id := range VoiceTypeRow {
			name := VoiceTypeMap[id]

			text += fmt.Sprintf("%v: %v\n", name, id)
		}

		text += fmt.Sprintf("\n以下のコマンドを送信することで、あなたのメッセとステータスを読み上げる音声を設定できます\n%vvoice <右側の数字>", GetPrefix())

		ctx.Reply(text, false)
	} else {
		voiceTypeNum, err := strconv.Atoi(ctx.Arguments[0])

		if err != nil {
			slog.Error(err.Error())
			return
		}

		if voiceType, ok := VoiceTypeMap[voiceTypeNum]; ok {
			VoiceSettings[ctx.Message.Author.ID] = VoiceSetting{VoiceType: voiceTypeNum}

			SaveVoiceSettings()

			ctx.Reply(fmt.Sprintf("あなたのメッセとステータスを読み上げる音声を「%v」に変更しました", voiceType), false)
		}
	}
}

func HandleBlockLists(ctx *disgolf.MessageCtx) {
	if len(ctx.Arguments) < 2 {
		text := "ブラックリスト\n\n"

		for _, id := range BlockList {
			user, _ := ctx.User(id)

			if user.ID != "" {
				text += fmt.Sprintf("%v: %v#%v\n", user.ID, user.Username, user.Discriminator)
			} else {
				text += fmt.Sprintf("%v: Unknown\n", user.ID)
			}
		}

		ctx.Reply(text, false)
	} else {
		mode := ctx.Arguments[0]

		switch mode {
		case "add":
			{
				text := ctx.Arguments[1]

				_, err := strconv.Atoi(text)

				if err != nil {
					slog.Error(err.Error())
					return
				}

				if slices.Contains(BlockList, text) {
					ctx.Reply(fmt.Sprintf("%vは既にブラックリストに追加されています", text), false)
				} else {
					BlockList = append(BlockList, text)

					SaveBlockLists()

					ctx.Reply(fmt.Sprintf("%vをブラックリストに追加しました", text), false)
				}
			}
		case "remove":
		case "delete":
		case "del":
			{
				text := ctx.Arguments[1]

				_, err := strconv.Atoi(text)

				if err != nil {
					slog.Error(err.Error())
					return
				}

				if slices.Contains(BlockList, text) {
					for i, id := range BlockList {
						if id == text {
							BlockList = append(BlockList[:i], BlockList[i+1:]...)
							break
						}
					}

					SaveBlockLists()

					ctx.Reply(fmt.Sprintf("%vをブラックリストから削除しました", text), false)
				} else {
					ctx.Reply(fmt.Sprintf("%vは元々ブラックリストにいません", text), false)
				}
			}

		}
	}
}

func HandleMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	/*if m.GuildID == "1006903965180510330" {
		content, err := m.ContentWithMoreMentionsReplaced(s)

		if err == nil {
			channel, err := s.Channel(m.ChannelID)

			if err == nil {
				var files []*discordgo.File

				for _, attach := range m.Message.Attachments {
					res, err := http.Get(attach.URL)

					if err == nil {
						file := &discordgo.File{Name: attach.Filename, ContentType: attach.ContentType, Reader: res.Body}
						files = append(files, file)
					}
				}

				s.WebhookExecute("1121370697248288839", "7YBvLYen1Io2hQhBwjRHrq_uh7Le7VPB1PuhiYGBg58QiogOrwOAU9nMSY0zTbAUz5-o", false, &discordgo.WebhookParams{Content: content, Username: fmt.Sprintf("#%v | %v#%v (%v)", channel.Name, m.Author.Username, m.Author.Discriminator, m.Author.ID), AvatarURL: m.Author.AvatarURL(""), Files: files, AllowedMentions: &discordgo.MessageAllowedMentions{}})
			}
		}
	}

	if m.ChannelID == "1089542986771218454" {
		fmt.Println(m.Content)
		fmt.Println(BotPrefixesPattern.FindStringIndex(m.Content))
	}*/

	//fmt.Println(m.Content)

	if m.Author.ID == s.State.User.ID || m.Author.Bot || strings.HasPrefix(m.Content, GetPrefix()) || strings.HasPrefix(m.Content, ";") || len(BotPrefixesPattern.FindStringIndex(m.Content)) != 0 || slices.Contains(BlockList, m.Author.ID) {
		return
	}

	voiceManager, err := GetManager(m.Message.GuildID)

	if err == nil {
		if voiceManager.Channel.Text.ID != m.ChannelID {
			return
		}

		if m.Content == "s" {
			voiceManager.Skip()
			return
		}

		if strings.HasPrefix(m.Content, ";") {
			return
		}

		content, err := m.ContentWithMoreMentionsReplaced(s)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		attachments := 0

		if len(m.Attachments) != 0 {
			attachments += len(m.Attachments)
		}

		attachments += len(DiscordFilePattern.FindAllStringIndex(content, -1))

		if attachments != 0 {
			//m.Attachments[0].ContentType
			content = fmt.Sprintf("%vつの添付ファイル。", attachments) + content
		}

		if m.ReferencedMessage != nil {
			content = "リプライ、" + content
		}

		content = EmojiPattern.ReplaceAllString(content, "$1")

		content = DiscordInvitePattern.ReplaceAllStringFunc(content, func(inviteCode string) string {
			texts := strings.Split(inviteCode, "/")

			guildName, err := GetInvite(texts[len(texts)-1])

			if err != nil {
				return "招待リンク、"
			}

			return fmt.Sprintf("招待リンク、%s、", guildName)
		})

		content = LinkPattern.ReplaceAllString(content, "$1のリンク")

		content = ExcludePattern.ReplaceAllString(content, "")

		content = CodeBlockPattern.ReplaceAllString(content, "コード省略")

		content = DuplicatedPattern.ReplaceAllString(content, "$1$1$1")

		fmt.Println(content)

		if utf8.RuneCountInString(content) < 100 {
			go voiceManager.Play([]string{content}, m.Author)
			return
		}

		var contents []string

		for {
			if utf8.RuneCountInString(content) >= 100 {
				splited := content[:99]

				contents = append(contents, splited)

				content = content[99:]
			} else {
				contents = append(contents, content)
				break
			}
		}

		go voiceManager.Play(contents, m.Author)
	}
}

func HandleVoiceUpdates(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	if m.Member.User.ID == s.State.User.ID || m.Member.User.Bot || slices.Contains(BlockList, m.UserID) {
		return
	}

	voiceManager, err := GetManager(m.GuildID)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	if m.BeforeUpdate == nil {
		m.BeforeUpdate = &discordgo.VoiceState{}
	}

	beforeID := m.BeforeUpdate.ChannelID
	nowID := m.ChannelID

	before, err := s.Channel(beforeID)

	if err != nil {
		before = &discordgo.Channel{}
	}

	now, err := s.Channel(nowID)

	if err != nil {
		now = &discordgo.Channel{}
	}

	name := m.Member.User.Username

	if m.Member.Nick != "" {
		name = m.Member.Nick
	}

	text := name

	if before.ID == voiceManager.Channel.Voice.ID || now.ID == voiceManager.Channel.Voice.ID {
		if now.ID != "" && !m.BeforeUpdate.SelfStream && m.SelfStream {
			text += "、が画面共有を開始しました"
		}

		if before.ID == "" && now.ID != "" {
			text += "、が入室しました"
		} else if before.ID != "" && now.ID == "" {
			text += "、が退出しました"
		} else if before.ID != "" && now.ID != "" {
			if now.ID == voiceManager.Channel.Voice.ID && before.ID != voiceManager.Channel.Voice.ID {
				text += "、が移動して来ました。"
			} else if now.ID != voiceManager.Channel.Voice.ID && before.ID == voiceManager.Channel.Voice.ID {
				text += "、が移動しました。"
			}
		}

	}

	if text != name {
		voiceManager.Play([]string{text}, m.Member.User)
	}
}
