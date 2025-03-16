package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func runPlay(ctx context.Context, cmd *cli.Command) error {
	file := cmd.Args().First()
	player := NewAutoPlayer(NewPipeWire(), NewPulseAudio())

	log.Debugf("playing file: %s", file)
	if err := player.Play(ctx, file, cmd.Float("volume")); err != nil {
		return fmt.Errorf("play: %w", err)
	}
	log.Debug("done playing")

	return nil
}
