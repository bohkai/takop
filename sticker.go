package main

import (
	"context"
	"strings"
	"log"

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

	parsed, err := Parse(m.Content)
	if err != nil {
		return
	}

	if parsed[0] != "s" {
		return
	}

	serchWord := strings.Join(parsed[1:], " ")
	ctx := context.Background()
	service, err := customsearch.NewService(ctx, option.WithAPIKey(st.GoogleConfig.Key))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Google API……なんで、死んだ？")
		return
	}
	serch := service.Cse.List()
	serch.Q(serchWord)
	serch.Cx(st.GoogleConfig.ID)
	serch.Num(1)
	serch.SearchType("image")

	call, err := serch.Do()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "探せなかったッピ……")
		return
	}

	if len(call.Items) == 0 {
		s.ChannelMessageSend(m.ChannelID, "探せなかったッピ……")
		return
	}

	go func ()  {
		if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
			log.Println(err)
			return;
		}
	}()

	go func ()  {
		user := GetUser(s, m.Message)
		link := call.Items[0].Link
		user.name = user.name + " (" + serchWord + ")"
		if err := user.AsSend(s, m.Message, link); err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "なかったッピ……")
			return;
		}
	}()
}
