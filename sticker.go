package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type sticker struct {
	*GoogleConfig
}

func NewSticker(config *GoogleConfig) *sticker {
	return &sticker{
		config,
	}
}

func (st *sticker) Serch(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, options, err := Parse(m.Content)
	if err != nil {
		return
	}
	if parsed[0] != "s" {
		return
	}

	imageIndex := 0
	if options != nil {
		i, err := strconv.Atoi(*options)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "オプションが不正だっピ……")
			return
		}
		imageIndex = i

		if imageIndex < 0 {
			s.ChannelMessageSend(m.ChannelID, "オプションの値が小さすぎるッピ!")
			return
		}

		if imageIndex >= 10 {
			s.ChannelMessageSend(m.ChannelID, "オプションの値が大きすぎるッピ!")
			return
		}
	}

	searchWord := strings.Join(parsed[1:], " ")
	ctx := context.Background()
	service, err := customsearch.NewService(ctx, option.WithAPIKey(st.GoogleConfig.Key))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Google API……なんで、死んだ？")
		return
	}
	search := service.Cse.List()
	search.Q(searchWord)
	search.Cx(st.GoogleConfig.ID)
	search.Num(10)
	search.SearchType("image")

	call, err := search.Do()
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "探せなかったッピ……")
		return
	}

	if len(call.Items) == 0 {
		s.ChannelMessageSend(m.ChannelID, "探せなかったッピ……")
		return
	}

	if imageIndex >= len(call.Items) {
		s.ChannelMessageSend(m.ChannelID, "これ以上画像はないッピ！我慢するッピ！")
		return
	}

	go func() {
		if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
			log.Println(err)
			return
		}
	}()

	go func() {
		link := call.Items[imageIndex].Link
		user := GetUser(s, m.Message)
		user.name = user.name + " (" + searchWord + ")"
		if err := user.AsSend(s, m.Message, link); err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "なかったッピ……")
			return
		}
	}()
}
