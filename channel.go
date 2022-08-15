package main

import (
	"context"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Channel struct {
	radiko *radiko
	se     *SE
	cancel context.CancelFunc
}

func NewChannel() (*Channel, error) {
	radiko, err := NewRadiko()
	se := NewSE()
	if err != nil {
		return nil, err
	}
	return &Channel{
		radiko,
		se,
		nil,
	}, nil
}

func (c *Channel) Join(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, _, err := Parse(m.Content)
	if err != nil {
		return
	}

	if len(parsed) != 1 || parsed[0] != "join" {
		return
	}

	_, err = c.ChannelVoiceJoin(s, m.Message)
	if err != nil {
		return
	}
}

func (c *Channel) Leave(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, _, err := Parse(m.Content)
	if err != nil {
		return
	}

	if len(parsed) != 1 || parsed[0] != "dis" {
		return
	}

	s.VoiceConnections[m.GuildID].Disconnect()
	c.Stop()
}

func (c *Channel) Play(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, _, err := Parse(m.Content)
	if err != nil {
		log.Println(err)
		return
	}

	if len(parsed) != 2 || parsed[0] != "play" {
		return
	}

	v, err := c.ChannelVoiceJoin(s, m.Message)
	if err != nil {
		log.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "チャンネルに入れないッピ、ちゃんとお話しするッピ ")
		return
	}

	c.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	err = c.radiko.RadikoPlay(s, m.Message, v, ctx, parsed[1])
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "なんで死んだ？\n"+err.Error())
	}
}

func (c *Channel) PlaySE(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	url := ""
	if len(m.Message.Attachments) == 1 {
		urls := m.Message.Attachments[0].URL
		if strings.HasPrefix(urls, "https://") &&
			(strings.HasSuffix(urls, ".mp3") || strings.HasSuffix(urls, ".wav")) {
			url = urls
		}
	}

	if m.Content != "" {
		urls := m.Content
		if strings.HasPrefix(urls, "https://") &&
			(strings.HasSuffix(urls, ".mp3") || strings.HasSuffix(urls, ".wav")) {
			url = urls
		}
	}

	if url == "" {
		return
	}

	v, err := c.ChannelVoiceJoin(s, m.Message)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "チャンネルに入れないッピ、ちゃんとお話しするッピ ")
		return
	}

	c.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.se.url = url
	c.se.Play(s, m.Message, v, ctx, url)
}

func (c *Channel) ReplaySE(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, _, err := Parse(m.Content)
	if err != nil {
		return
	}

	if len(parsed) != 1 || parsed[0] != "r" {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	v, err := c.ChannelVoiceJoin(s, m.Message)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "チャンネルに入れないッピ、ちゃんとお話しするッピ ")
		return
	}

	c.se.RePlay(s, m.Message, v, ctx)
}

func (c *Channel) List(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, _, err := Parse(m.Content)
	if err != nil {
		return
	}

	if len(parsed) != 1 || parsed[0] != "list" {
		return
	}

	c.radiko.RadikoList(s, m.Message)
}

func (c *Channel) ChannelVoiceJoin(s *discordgo.Session, m *discordgo.Message) (*discordgo.VoiceConnection, error) {
	vs, err := s.State.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return nil, err
	}
	return s.ChannelVoiceJoin(m.GuildID, vs.ChannelID, false, false)
}

func (c *Channel) StopSound(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, _, err := Parse(m.Content)
	if err != nil {
		return
	}

	if len(parsed) != 1 || parsed[0] != "stop" {
		return
	}

	c.Stop()
}

func (c *Channel) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
}
