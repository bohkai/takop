package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	goradiko "github.com/yyoshiki41/go-radiko"
)

type radiko struct {
	client *goradiko.Client
}

func NewRadiko() (*radiko, error) {
	client, err := goradiko.New("")
	if err != nil {
		return nil, err
	}

	return &radiko{
		client,
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
	if channel == "" {
		return errors.New("idを入れると幸せになれるッピ！")
	}

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
		return errors.New("配信中のストリームが見つからないッピ！")
	}

	_, err = r.client.AuthorizeToken(context.Background())
	if err != nil {
		return errors.New("authorize token error")
	}
	token := r.client.AuthToken()

	ffmpegCmd, err := NewFfmpeg()
	if err != nil {
		return err
	}
	ffmpegArgs := []string{
		"-headers", "X-Radiko-Authtoken: " + token,
		"-i", streamURL,
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
	}

	ffmpegCmd.SetArgs(ffmpegArgs...)
	ffmpegout, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		return err
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	go func() {
		err = ffmpegCmd.Run("pipe:1")
		if err != nil {
			log.Println("ffmpeg error:" + err.Error())
			s.ChannelMessageSend(m.ChannelID, "ffmpegが死んでるッピ！")
		}
	}()

	go func(ctx context.Context) {
		<-ctx.Done()
		log.Println("ffmpeg done")
		err = ffmpegCmd.Kill()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "完膚なきまでにffmpegを壊さなきゃ")
			log.Println("ffmpeg kill error:" + err.Error())
			return
		}
	}(ctx)

	go func(ctx context.Context) {
		v.Speaking(true)
		send := make(chan []int16, 2)
		defer close(send)
		defer v.Speaking(false)

		go func() {
			dgvoice.SendPCM(v, send)
		}()

		for {
			audiobuf := make([]int16, 960*2)
			if err := binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf); err != nil {
				s.ChannelMessageSend(m.ChannelID, "binaryが わ わかんないッピ……")
				log.Println("binary.Read error:" + err.Error())
				return
			}
			select {
			case send <- audiobuf:
				continue
			case <-ctx.Done():
				log.Println("ctx done")
				return
			}
		}
	}(ctx)

	return nil
}
