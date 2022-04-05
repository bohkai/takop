package main

import (
	"context"
	"os/exec"
)

type ffmpeg struct {
	*exec.Cmd
	isPlay chan bool
}

func NewFfmpeg(ctx context.Context) (*ffmpeg, error) {
	cmdPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}

	return &ffmpeg{
		exec.CommandContext(
			ctx,
			cmdPath,
		),
		make(chan bool)}, nil
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

func (f *ffmpeg) Stop() {
	f.isPlay <- true
	close(f.isPlay)
}
