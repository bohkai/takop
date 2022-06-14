package main

import (
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"context"
	"bufio"
	"log"
)

type SE struct {
	url string
}

func NewSE() *SE {
 return &SE{}
}

func (e *SE) Play(s *discordgo.Session, m *discordgo.Message, v *discordgo.VoiceConnection, ctx context.Context, url string) error {
	e.setURL(url)
	ffmpegCmd, err := NewFfmpeg()
	if err != nil {
		return err
	}

	ffmpegArgs := []string{
		"-i", url,
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
	err = ffmpegCmd.Start("pipe:1")
	if err != nil {
		log.Println("ffmpeg error:" + err.Error())
		return err
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		log.Println("ffmpeg done")
		err = ffmpegCmd.Kill()
		if err != nil {
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

		err := ffmpegCmd.Play(ffmpegbuf, send, ctx)
		if err != nil {
			return
		}
	}(ctx)

	return nil
}

func (e *SE) RePlay(s *discordgo.Session, m *discordgo.Message, v *discordgo.VoiceConnection, ctx context.Context) {
	if e.url == "" {
		return
	}
	e.Play(s, m, v, ctx, e.url)
}

func (e *SE) setURL(name string) {
	e.url = name
}
