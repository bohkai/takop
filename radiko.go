package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/yyoshiki41/go-radiko"
)

func RadioList(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	client, err := radiko.New("")
	if err != nil {
		return
	}

	stations, err := client.GetNowPrograms(context.Background())
	if err != nil {
		return
	}

	message := ""
	for _, station := range stations {
		message = message + fmt.Sprintf("%s | %30.30s | %-30.20s\n", station.ID, station.Name, station.Scd.Progs.Progs[0].Title)
	}
	s.ChannelMessageSend(m.ChannelID, message)
	context.Background().Done()
}

func RadioPlay(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parsed, err := Parse(m.Content)
	if err != nil {
		return
	}

	if len(parsed) != 2 || parsed[0] != "play" {
		return
	}

	ctx := context.Background()
	client, err := radiko.New("")
	if err != nil {
		return
	}
	_, err = client.AuthorizeToken(ctx)
	if err != nil {
		return
	}

	items, err := radiko.GetStreamSmhMultiURL(parsed[1])
	if err != nil {
		return
	}
	var streamURL string
	for _, item := range items {
		if !item.Areafree {
			streamURL = item.PlaylistCreateURL
			break
		}
	}

	if streamURL == "" {
		return
	}

	ffmpegCmd, err := NewFfmpeg(ctx, streamURL)
	if err != nil {
		log.Println(err)
		return
	}

	ffmpegArgs := []string{
		"-headers", "X-Radiko-Authtoken: " + client.AuthToken(),
		"-i", streamURL,
		"-y",
		"-vn",
		"-acodec",
		"copy",
	}
	ffmpegCmd.setArgs(ffmpegArgs...)
	err = ffmpegCmd.Run("./test.m4a")
	if err != nil {
		log.Println("ffmpeg error:" + err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Playing")
}