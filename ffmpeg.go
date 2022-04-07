package main

import (
	"context"
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
		),}, nil
}

func (f *ffmpeg) SetArgs(args ...string) {
	f.Args = append(f.Args, args...)
}

func (f *ffmpeg) Run(output string) error {
	f.SetArgs(output)
	return f.Cmd.Run()
}

func (f *ffmpeg) Kill() error {
	return f.Cmd.Process.Kill()
}
