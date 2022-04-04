package main

import (
	"context"
	"os/exec"
)

type ffmpeg struct {
	*exec.Cmd
}

func NewFfmpeg(ctx context.Context, url string) (*ffmpeg, error) {
	cmdPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}

	return &ffmpeg{exec.CommandContext(
		ctx,
		cmdPath,
	)}, nil
}

func (f *ffmpeg) setArgs(args ...string) {
	f.Args = append(f.Args, args...)
}

func (f *ffmpeg) Run(output string) error {
	f.setArgs(output)
	return f.Cmd.Run()
}
