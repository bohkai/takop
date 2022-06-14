package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"os/exec"
)

type ffmpeg struct {
	*exec.Cmd
}

func NewFfmpeg() (*ffmpeg, error) {
	cmdPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}

	return &ffmpeg{
		exec.CommandContext(
			context.Background(),
			cmdPath,
		),
		}, nil
}

func (f *ffmpeg) SetArgs(args ...string) {
	f.Args = append(f.Args, args...)
}

func (f *ffmpeg) Start(output string) error {
	f.SetArgs(output)
	return f.Cmd.Start()
}

func (f *ffmpeg) Kill() error {
	return f.Cmd.Process.Kill()
}

func (f *ffmpeg) Play(buf *bufio.Reader, send chan[]int16 , ctx context.Context) error {
	for {
		audiobuf := make([]int16, 960*2)
		if err := binary.Read(buf, binary.LittleEndian, &audiobuf); err != nil {
			return err
		}
		select {
		case send <- audiobuf:
			continue
		case <-ctx.Done():
			return nil
		}
	}
}
