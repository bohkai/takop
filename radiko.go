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
	IsRadioStop     chan bool
	ctx             context.Context
}

func NewRadiko(ctx context.Context) (*radiko, error) {
	client, err := goradiko.New("")
	if err != nil {
		return nil, err
	}

	_, err = client.AuthorizeToken(ctx)
	if err != nil {
		return nil, err
	}

	return &radiko{
		client,
		client.AuthToken(),
		make(chan bool),
		make(chan bool),
		ctx}, nil
}

func (r *radiko) RadikoList(s *discordgo.Session, m *discordgo.MessageCreate) {
	stations, err := r.client.GetNowPrograms(r.ctx)
	if err != nil {
		return
	}
	message := ""
	for _, station := range stations {
		message = message + fmt.Sprintf("%s | %30.30s | %-30.20s\n", station.ID, station.Name, station.Scd.Progs.Progs[0].Title)
	}
	s.ChannelMessageSend(m.ChannelID, message)
}

func (r *radiko) RadikoPlay(s *discordgo.Session, m *discordgo.MessageCreate, v *discordgo.VoiceConnection, channel string, ffmpegCmd *ffmpeg) error {
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

	ffmpegArgs := []string{
		"-headers", "X-Radiko-Authtoken: " + r.token,
		"-i", streamURL,
		"-f",
		"segment", "-segment_time", "10",
		"-y",
		"-vn",
		"-acodec",
		"copy",
	}
	ffmpegCmd.SetArgs(ffmpegArgs...)

	go func() {
		err = ffmpegCmd.Start("./audio/out-%d.m4a")
		if err != nil {
			log.Println("ffmpeg error:" + err.Error())
			return
		}
		for {
			select {
			case <- ffmpegCmd.isPlay: ffmpegCmd.Process.Kill()
			default: continue
			}
		}
	}()

	s.ChannelMessageSend(m.ChannelID, "バッファリング中... 30秒お待ち下さい")
	time.Sleep(time.Second * 10)
	s.ChannelMessageSend(m.ChannelID, "再生します")

	go func() {
		number := 0
		for {
			path := fmt.Sprintf("./audio/out-%d.m4a", number)
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				s.ChannelMessageSend(m.ChannelID, "バッファリング中... 30秒お待ち下さい")
				time.Sleep(time.Second * 30)
				continue
			}
			dgvoice.PlayAudioFile(v, path, r.IsVoicePlayStop)
			os.Remove(path)
			number++

			select {
			case <-r.IsRadioStop:
				return
			default:
			}
		}
	}()

	return nil
}

func (r *radiko) RadikoStop() {
	r.IsVoicePlayStop <- true
	r.IsRadioStop <- true
	close(r.IsVoicePlayStop)
}
