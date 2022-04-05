package main

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Channel struct {
	ffmpeg *ffmpeg
	radiko *radiko
	ctx context.Context
}

func NewChannel(ctx context.Context) (*Channel, error) {
	radiko, err := NewRadiko(ctx)
	if err != nil {
		return nil, err
	}

	return &Channel{
		nil,
		radiko,
		ctx,
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

	v, err := c.ChannelVoiceJoin(s, m)
	if err != nil {
		log.Println(err)
		return
	}

	ffmpeg, err := NewFfmpeg(c.ctx)
	if err != nil {
		log.Println(err)
		return
	}
	c.ffmpeg = ffmpeg
	err = c.radiko.RadikoPlay(s, m, v, parsed[1], c.ffmpeg)
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
