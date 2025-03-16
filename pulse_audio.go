package main

import (
	"context"
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/log"
)

type PulseAudio struct{}

func NewPulseAudio() *PulseAudio {
	return &PulseAudio{}
}

func (p *PulseAudio) Name() string {
	return "PulseAudio"
}

func (p *PulseAudio) Supported() bool {
	_, err := exec.LookPath("paplay")
	return err == nil
}

func (p *PulseAudio) Play(ctx context.Context, file string, volume float64) error {
	cmd := exec.CommandContext(ctx,
		"paplay",
		"--volume", strconv.FormatInt(int64(volume*65536), 10),
		file,
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *PulseAudio) PlayInBackground(file string, volume float64) error {
	cmd := exec.Command(
		"paplay",
		"--volume", strconv.FormatInt(int64(volume*65536), 10),
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
