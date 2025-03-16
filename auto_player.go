package main

import (
	"context"
	"errors"

	"github.com/charmbracelet/log"
)

type AutoPlayer struct {
	players []Player
}

func NewAutoPlayer(players ...Player) *AutoPlayer {
	return &AutoPlayer{
		players: players,
	}
}

func (a *AutoPlayer) Name() string {
	for _, player := range a.players {
		if player.Supported() {
			return player.Name() + " (AutoPlayer)"
		}
	}
	return "none (AutoPlayer)"
}

func (a *AutoPlayer) Supported() bool {
	for _, player := range a.players {
		if player.Supported() {
			return true
		}
	}
	return false
}

func (a *AutoPlayer) Play(ctx context.Context, file string, volume float64) error {
	for _, player := range a.players {
		if player.Supported() {
			log.Debugf("playing with %s", player.Name())
			return player.Play(ctx, file, volume)
		}
	}
	return errors.New("no supported player")
}

func (a *AutoPlayer) PlayInBackground(file string, volume float64) error {
	for _, player := range a.players {
		if player.Supported() {
			log.Debugf("playing in background with %s", player.Name())
			return player.PlayInBackground(file, volume)
		}
	}
	return errors.New("no supported player")
}
