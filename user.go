package main

import (
	"github.com/bwmarrin/discordgo"
)

type User struct {
	name string
	avatar string
}

func (u *User) AsSend(s *discordgo.Session, m *discordgo.Message, content string) error {
	st, err := s.WebhookCreate(m.ChannelID, u.name, u.avatar)
	if err != nil {
		return err
	}

	body := &discordgo.WebhookParams{
		Content: content,
		Username: u.name,
		AvatarURL: u.avatar,
	}
	if _, err := s.WebhookExecute(st.ID, st.Token, false, body); err != nil {
		return err
	}
	s.WebhookDelete(st.ID)
	return nil
}

func GetUser(s *discordgo.Session, m *discordgo.Message) (*User) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		return &User{
			name: m.Author.Username,
			avatar: m.Author.AvatarURL(""),
		}
	}
	name := member.Nick
	if name == "" {
		name = member.User.Username
	}
	return &User{
		name: name,
		avatar: member.User.AvatarURL(""),
	}
}
