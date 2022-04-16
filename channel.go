package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
)

type Channel struct {
	radiko *radiko
	cancel context.CancelFunc
}

func NewChannel() (*Channel, error) {
	radiko, err := NewRadiko()
	if err != nil {
		return nil, err
	}
	return &Channel{
		radiko,
		nil,
	}, nil
}

func (c *Channel) Join(s *discordgo.Session, m *discordgo.Message) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, err := Parse(m.Content)
	if err != nil {
		return
	}

	if len(parsed) != 1 || parsed[0] != "join" {
		return
	}

	_, err = c.ChannelVoiceJoin(s, m)
	if err != nil {
		return
	}
}

func (c *Channel) Leave(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, err := Parse(m.Content)
	if err != nil {
		return
	}

	if len(parsed) != 1 || parsed[0] != "dis" {
		return
	}

	s.VoiceConnections[m.GuildID].Disconnect()
	c.Stop()
}

func (c *Channel) Play(s *discordgo.Session, m *discordgo.Message) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, err := Parse(m.Content)
	if err != nil {
		log.Println(err)
		return
	}

	if len(parsed) != 2 || parsed[0] != "play" {
		return
	}

	c.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	v, err := c.ChannelVoiceJoin(s, m)
	if err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "チャンネルに入れないッピ、ちゃんとお話しするッピ ")
		return
	}

	err = c.radiko.RadikoPlay(s, m, v, ctx, parsed[1])
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "なんで死んだ？\n" + err.Error())
		return
	}
}

func (c *Channel) List(s *discordgo.Session, m *discordgo.Message) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, err := Parse(m.Content)
	if err != nil {
		return
	}

	if len(parsed) != 1 || parsed[0] != "list" {
		return
	}

	c.radiko.RadikoList(s, m)
}

func (c *Channel) ChannelVoiceJoin(s *discordgo.Session, m *discordgo.Message) (*discordgo.VoiceConnection, error) {
	vs, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return nil, err
	}
	return s.ChannelVoiceJoin(m.GuildID, vs.ChannelID, false, false)
}

func (c *Channel) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}
