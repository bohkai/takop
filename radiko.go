package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	goradiko "github.com/yyoshiki41/go-radiko"
)

type radiko struct {
	client          *goradiko.Client
	token           string
	IsVoicePlayStop chan bool
}

func NewRadiko() (*radiko, error) {
	client, err := goradiko.New("")
	if err != nil {
		return nil, err
	}

	_, err = client.AuthorizeToken(context.Background())
	if err != nil {
		return nil, err
	}

	return &radiko{
		client,
		client.AuthToken(),
		make(chan bool),
	}, nil
}

func (r *radiko) RadikoList(s *discordgo.Session, m *discordgo.MessageCreate) {
	stations, err := r.client.GetNowPrograms(context.Background())
	if err != nil {
		return
	}
	message := ""
	for _, station := range stations {
		message = message + fmt.Sprintf("%s | %30.30s | %-30.20s\n", station.ID, station.Name, station.Scd.Progs.Progs[0].Title)
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

func (r *radiko) RadikoPlay(s *discordgo.Session, m *discordgo.MessageCreate, v *discordgo.VoiceConnection, ctx context.Context, channel string) error {
	items, err := goradiko.GetStreamSmhMultiURL(channel)
	if err != nil {
		return err
	}

	var streamURL string
	for _, item := range items {
		if !item.Areafree {
			streamURL = item.PlaylistCreateURL
			break
		}
	}

	if streamURL == "" {
		return errors.New("no stream URL")
	}

	audioPath := "./audio"
	os.RemoveAll(audioPath)
	_, err = os.Stat(audioPath)
	if os.IsNotExist(err) {
		os.Mkdir(audioPath, 0777)
	}

	ffmpegCmd, err := NewFfmpeg()
	if err != nil {
		return err
	}
	ffmpegArgs := []string{
		"-headers", "X-Radiko-Authtoken: " + r.token,
		"-i", streamURL,
		"-f",
		"segment", "-segment_time", "30",
		"-y",
		"-vn",
		"-acodec",
		"copy",
	}
	ffmpegCmd.SetArgs(ffmpegArgs...)

	go func(ctx context.Context) {
		err = ffmpegCmd.Start("./audio/out-%d.m4a")
		if err != nil {
			log.Println("ffmpeg error:" + err.Error())
			return
		}
		<-ctx.Done()
		log.Println("ffmpeg done")
		err = ffmpegCmd.Kill()
		if err != nil {
			log.Println("ffmpeg kill error:" + err.Error())
			return
		}
	}(ctx)

	go func(ctx context.Context) {
		number := 0
		path := ""

		for {
			path = fmt.Sprintf("./audio/out-%d.m4a", number)
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				t := time.NewTicker(time.Second * 30)
				s.ChannelMessageSend(m.ChannelID, "バッファリング中... 30秒お待ち下さい")
				L: for {
					select {
					case <-t.C:
						if v.Ready {
							s.ChannelMessageSend(m.ChannelID, "再生開始")
							break L
						}
					case <-ctx.Done():
						log.Println("radiko done")
						return
					}
				}
				continue
			}
			dgvoice.PlayAudioFile(v, path, r.IsVoicePlayStop)
			number++

			select {
			case <-ctx.Done():
				log.Println("radiko done")
				return
			case <- r.IsVoicePlayStop:
			default:
				os.Remove(path)
			}
		}
	}(ctx)

	return nil
}
