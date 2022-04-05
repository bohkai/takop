package main

import (
	"log"
	"github.com/bwmarrin/discordgo"
	"context"
)

type Channel struct {
	radiko *radiko
	cancel  context.CancelFunc
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

func (c *Channel)Join(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	for _, g := range s.State.Guilds {
		for _, vs := range g.VoiceStates {
			if m.Author.ID != vs.UserID {
				continue
			}
			s.VoiceConnections[vs.GuildID].Disconnect()
			c.Stop()
		}
	}
}

func (c *Channel) Play(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		log.Println(err)
		return
	}

	err = c.radiko.RadikoPlay(s, m, v, ctx, parsed[1])
	if err != nil {
		log.Println(err)
		return
	}
}

func (c *Channel) List(s *discordgo.Session, m *discordgo.MessageCreate) {
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

func (c *Channel) ChannelVoiceJoin(s *discordgo.Session, m *discordgo.MessageCreate) (*discordgo.VoiceConnection, error) {
	for _, g := range s.State.Guilds {
		for _, vs := range g.VoiceStates {
			if m.Author.ID != vs.UserID {
				continue
			}
			return s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, false)
		}
	}

	return nil, nil
}

func (c *Channel) Stop() {
	if c.cancel != nil {
		c.cancel()
	}

	select {
	case c.radiko.IsVoicePlayStop <- true:
		close(c.radiko.IsVoicePlayStop)
		c.radiko.IsVoicePlayStop = make(chan bool)
	default:
	}
}