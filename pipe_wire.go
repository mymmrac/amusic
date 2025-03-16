package main

import (
	"context"
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/log"
)

type PipeWire struct{}

func NewPipeWire() *PipeWire {
	return &PipeWire{}
}

func (p *PipeWire) Name() string {
	return "PipeWire"
}

func (p *PipeWire) Supported() bool {
	_, err := exec.LookPath("pw-play")
	return err == nil
}

func (p *PipeWire) Play(ctx context.Context, file string, volume float64) error {
	cmd := exec.CommandContext(ctx,
		"pw-play",
		"--volume", strconv.FormatFloat(volume, 'f', -1, 64),
		file,
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PipeWire) PlayInBackground(file string, volume float64) error {
	cmd := exec.Command(
		"pw-play",
		"--volume", strconv.FormatFloat(volume, 'f', -1, 64),
		file,
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Errorf("play in background: %s", err)
			return
		}
		log.Debug("done playing in background")
	}()

	return nil
}
