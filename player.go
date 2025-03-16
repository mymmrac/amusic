package main

import "context"

type Player interface {
	Name() string
	Supported() bool
	Play(ctx context.Context, file string, volume float64) error
	PlayInBackground(file string, volume float64) error
}
